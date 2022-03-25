let ws = undefined;
let game_id = undefined;
let sudoku = undefined;

document.addEventListener('DOMContentLoaded', () => {
    sudoku = document.querySelector('#sudoku');

    // Creating board in table element.
    for (let row = 0; row < 9; row++) {
        let tr = document.createElement('tr');
        for (let col = 0; col < 9; col++) {
            let td = document.createElement('td');
            td.id = String.fromCharCode('a'.charCodeAt(0)+row)+(col+1);
            tr.appendChild(td);
        }
        sudoku.appendChild(tr);
    }

    // Creating event handlers for all cells.
    sudoku.querySelectorAll('tr td').forEach((td) => {
        td.addEventListener('mouseup', function(e) {
            setActive(td);
        });
    });

    let placeDigitInActive = (digit) => {
        if (sudoku.classList.contains('win')) return;
        let td = sudoku.querySelector('td.active');
        if (td || td.classList.contains('hint')) return;
        let oldDigit = td.textContent;
        td.textContent = digit;
        if (oldDigit === digit) apiMakeStep();
    };

    let keyboard = document.querySelector('#mobile-keyboard');
    if (keyboard !== undefined) {
        let createButton = (id, label, event) => {
            let button = document.createElement('span');
            button.classList.add('keyboard-button');
            button.id = id;
            button.textContent = label;
            button.addEventListener('mouseup', event);
            keyboard.appendChild(button);
        }
        for (let digit = 1; digit <= 9; digit++) {
            createButton('button'+digit, digit, (e) => {
                placeDigitInActive(digit);
            });
        }
        createButton('button'+digit, '⨯', (e) => {
            placeDigitInActive('0');
        });
    }

    // Creating keyup and digit input handlers on document.
    document.addEventListener('keydown', (e) => {
        if (e.defaultPrevented) {
            return;
        }
        let td = document.querySelector('#sudoku tr td.active');
        switch (e.code) {
            case 'ArrowUp':    setActive(td, 'up');    break;
            case 'ArrowRight': setActive(td, 'right'); break;
            case 'ArrowDown':  setActive(td, 'down');  break;
            case 'ArrowLeft':  setActive(td, 'left');  break;
            case 'Digit0':
            case 'Numpad0':
            case 'Space':
            case 'Backspace':
                placeDigitInActive('0');
        }
        if ('1' <= e.key && e.key <= '9') {
            placeDigitInActive(e.key);
        }
    });

    sudoku.addEventListener('api_getPuzzle', (e) => {
        let body = e.detail.body;
        sudoku.querySelectorAll('tr').forEach((tr, row) => {
            tr.querySelectorAll('td').forEach((td, col) => {
                td.textContent = '';
                let d = body.puzzle[row*9+col];
                if ('1' <= d && d <= '9') {
                    td.textContent = d;
                    td.classList.add('hint');
                }
            });
        });
    });

    sudoku.addEventListener('api_makeStep', (e) => {
        let body = e.detail.body;
        sudoku.querySelectorAll('tr td').forEach((td) => {
            td.classList.remove('error');
        });
        if (body.win) {
            sudoku.classList.add('win');
            return;
        }
        if (!body.errors) {
            return;
        }
        body.errors = parsePoints(body.errors);
        sudoku.querySelectorAll('tr').forEach((tr, row) => {
            tr.querySelectorAll('td').forEach((td, col) => {
                body.errors.forEach((p) => {
                    if (p.row === row && p.col === col) {
                        td.classList.add('error');
                    }
                });
            });
        });
    });

    // websocket
    connectWs();
    // setInterval(()=>{
    //     if(!ws){
    //         return;
    //     }
    //     ws.send(JSON.stringify({method: 'health', echo: ''+Math.floor(Math.random() * 1e9)}));
    // }, 10000);

    sudoku.addEventListener('apiReady', () => {
        game_id = document.querySelector('#_game_id').textContent;
        wsApi('getPuzzle', {
            game_id: game_id,
        });
    }, {once: true});
}, false);

let setActive = (td, dir) => {
    if (!td) {
        td = sudoku.querySelectorAll('tr').item(9/2).querySelectorAll('td').item(9/2);
        dir = undefined;
        if (!td) return;
    }
    let tr = td.closest('tr');
    if (dir) {
        switch (dir) {
            case 'up':
                let prev = tr.previousElementSibling;
                if (!prev) return;
                td = prev.querySelectorAll('td').item(getIndex(td));
                break;
            case 'right':
                td = td.nextElementSibling; break;
            case 'down':
                let next = tr.nextElementSibling;
                if (!next) return;
                td = next.querySelectorAll('td').item(getIndex(td));
                break;
            case 'left':
                td = td.previousElementSibling; break;
        }
    }
    if (!td) return;
    let isAlreadyActive = td.classList.contains('active');
    document.querySelectorAll('#sudoku tr td.active').forEach((active) => {
        active.classList.remove('active');
    });
    if (!isAlreadyActive) td.classList.add('active');
}

function getIndex(node) {
    let index = 0;
    while (node = node.previousElementSibling) {
        index++;
    }
    return index;
}

let connectWs = () => {
    ws = new WebSocket('ws://'+location.host+'/ws');
    ws.onopen = (e) => {
        console.log('ws: open connection');
        sudoku.dispatchEvent(new CustomEvent('apiReady'));
    }
    ws.onclose = (e) => {
        console.log('ws: close connection');
        ws = undefined;
        // reconnect
        setTimeout(connectWs, 3000);
    }
    ws.onmessage = (e) => {
        console.log('ws: receive message:', e.data);
        let msg = JSON.parse(e.data);
        if (!msg.method) return;
        if (msg.error) {
            console.error('api', msg.method, 'error:', msg.error);
            return;
        }
        sudoku.dispatchEvent(new CustomEvent('api_'+msg.method, {detail: msg}));
    }
    ws.onerror = (e) => {
        console.error('ws: error '+e.code+':', e.reason, e);
        ws.close();
    }
}

let apiMakeStep = () => {
    let state = '';
    sudoku.querySelectorAll('tr td').forEach((td) => {
        let val = td.textContent;
        if (val === '') val = '.';
        state += val;
    });
    wsApi('makeStep', {
        game_id: game_id,
        state: state,
    })
}

let wsApi = (method, body) => {
    if (!body || !method) return;
    let msg = JSON.stringify({
        method: method,
        body: body,
    });
    console.log('ws: send message:', msg);
    ws.send(msg);
}

let parsePoints = (points) => {
    let out = [];
    points.forEach((p) => {
        out = out.concat([{
            row: p[0].charCodeAt(0)-'a'.charCodeAt(0),
            col: parseInt(p[1])-1,
        }]);
    });
    return out;
}
