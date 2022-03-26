let ws = undefined;
let game_id = undefined;
let sudoku = undefined;

document.addEventListener('DOMContentLoaded', () => {
    sudoku = document.querySelector('#sudoku');

    // Creating board in table element.
    for (let row = 0; row < 9; row++) {
        let elRow = document.createElement('div');
        elRow.classList.add('row');
        for (let col = 0; col < 9; col++) {
            let elCell = document.createElement('div');
            elCell.classList.add('cell');
            elCell.id = String.fromCharCode('a'.charCodeAt(0)+row)+(col+1); // TODO remove

            // digit place
            let elDigit = document.createElement('div');
            elDigit.classList.add('digit');
            elCell.appendChild(elDigit);

            // candidates place
            let elCands = document.createElement('div');
            elCands.classList.add('cands');
            for (let crow = 0; crow < 3; crow++) {
                let elCandsRow = document.createElement('div');
                for (let ccol = 0; ccol < 3; ccol++) {
                    let elCand = document.createElement('div');
                    let cval = crow*3+ccol+1;
                    elCand.classList.add('c'+cval);
                    elCand.textContent = ''+cval;
                    elCandsRow.appendChild(elCand);
                }
                elCands.appendChild(elCandsRow);
            }
            elCell.appendChild(elCands);

            elRow.appendChild(elCell);
        }
        sudoku.appendChild(elRow);
    }

    // Creating event handlers for all cells.
    sudoku.querySelectorAll('.cell').forEach((cell) => {
        cell.addEventListener('mouseup', function(e) {
            setActive(cell);
        });
    });

    let placeDigit = (elCell, digit, notMakeStep) => {
        if (sudoku.classList.contains('win')) return;
        if (!elCell || elCell.classList.contains('hint')) return;
        let elDigit = elCell.querySelector('.digit');
        let oldDigit = elDigit.textContent===''?'0':elDigit.textContent;
        elCell.classList.remove('is_digit', 'is_cands');
        if (digit === '0') {
            elDigit.textContent = '';
            elCell.classList.add('is_cands');
        } else {
            elDigit.textContent = digit;
            elCell.classList.add('is_digit');
        }
        if (!notMakeStep && oldDigit !== digit) apiMakeStep();
    };

    let placeDigitInActive = (digit, notMakeStep) => {
        if (sudoku.classList.contains('win')) return;
        placeDigit(sudoku.querySelector('.cell.active'), digit, notMakeStep);
    };

    let keyboard = document.querySelector('#keyboard');
    if (keyboard !== undefined) {
        let createButton = (id, label, event) => {
            let button = document.createElement('div');
            button.classList.add('keyboard-button');
            button.id = id;
            button.textContent = label;
            button.addEventListener('click', event);
            keyboard.appendChild(button);
        }
        createButton('buttonC', 'c', (e) => {
            console.log('TODO set candidates');
        });
        createButton('button0', '⨯', (e) => {
            placeDigitInActive('0');
        });
        for (let digit = 1; digit <= 9; digit++) {
            createButton('button'+digit, digit, (e) => {
                placeDigitInActive(''+digit);
            });
        }
    }

    // Creating keyup and digit input handlers on document.
    document.addEventListener('keydown', (e) => {
        if (e.defaultPrevented) {
            return;
        }
        let active = document.querySelector('#sudoku .cell.active');
        switch (e.code) {
            case 'ArrowUp':    setActive(active, 'up');    break;
            case 'ArrowRight': setActive(active, 'right'); break;
            case 'ArrowDown':  setActive(active, 'down');  break;
            case 'ArrowLeft':  setActive(active, 'left');  break;
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
        sudoku.querySelectorAll('.row').forEach((elRow, row) => {
            elRow.querySelectorAll('.cell').forEach((elCell, col) => {
                placeDigit(elCell, '0', true);
                let d = body.puzzle[row*9+col];
                if ('1' <= d && d <= '9') {
                    placeDigit(elCell, d, true);
                    elCell.classList.add('hint');
                }
            });
        });
    });

    sudoku.addEventListener('api_makeStep', (e) => {
        let body = e.detail.body;
        sudoku.querySelectorAll('.cell').forEach((cell) => {
            cell.classList.remove('error');
        });
        if (body.win) {
            sudoku.classList.add('win');
            return;
        }
        if (!body.errors) {
            return;
        }
        body.errors = parsePoints(body.errors);
        sudoku.querySelectorAll('.row').forEach((elRow, row) => {
            elRow.querySelectorAll('.cell').forEach((elCell, col) => {
                body.errors.forEach((p) => {
                    if (p.row === row && p.col === col) {
                        elCell.classList.add('error');
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

let setActive = (elCell, dir) => {
    if (!elCell) {
        elCell = sudoku.querySelectorAll('.row').item(9/2).querySelectorAll('.cell').item(9/2);
        dir = undefined;
        if (!elCell) return;
    }
    let elRow = elCell.closest('.row');
    if (dir) {
        switch (dir) {
            case 'up':
                let prev = elRow.previousElementSibling;
                if (!prev) return;
                elCell = prev.querySelectorAll('.cell').item(getIndex(elCell));
                break;
            case 'right':
                elCell = elCell.nextElementSibling; break;
            case 'down':
                let next = elRow.nextElementSibling;
                if (!next) return;
                elCell = next.querySelectorAll('.cell').item(getIndex(elCell));
                break;
            case 'left':
                elCell = elCell.previousElementSibling; break;
        }
    }
    if (!elCell) return;
    let isAlreadyActive = elCell.classList.contains('active');
    sudoku.querySelectorAll('.cell.active').forEach((active) => {
        active.classList.remove('active');
    });
    if (!isAlreadyActive) elCell.classList.add('active');
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
    sudoku.querySelectorAll('.cell .digit').forEach((cell) => {
        let val = cell.textContent;
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
