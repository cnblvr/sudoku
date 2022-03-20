let ws = undefined;

document.addEventListener('DOMContentLoaded', () => {
    let sudoku = document.querySelector('#sudoku');

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
    document.querySelectorAll('#sudoku tr td').forEach((td) => {
        td.addEventListener("mouseup", function(e) {
            setActive(td);
        });
    });

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
                if (td) td.textContent = '';
        }
        if (td && '1' <= e.key && e.key <= '9') {
            td.textContent = e.key;
        }
    });

    // websocket
    connectWs();
    // setInterval(()=>{
    //     if(!ws){
    //         return;
    //     }
    //     ws.send(JSON.stringify({method: 'health', echo: ''+Math.floor(Math.random() * 1e9)}));
    // }, 10000);
}, false);

let setActive = (td, dir) => {
    if (!td) {
        td = document.querySelectorAll('#sudoku tr').item(9/2).querySelectorAll('td').item(9/2);
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
    }
    ws.onclose = (e) => {
        console.log('ws: close connection');
        ws = undefined;
        // reconnect
        setTimeout(() => {
            connectWs();
        }, 3000);
    }
    ws.onmessage = (e) => {
        console.log('ws: new message: ', e.data);
    }
    ws.onerror = (e) => {
        console.error('ws: error('+e.code+'): ', e.reason, e);
        ws.close();
    }
}
