document.addEventListener('DOMContentLoaded', function() {
    let sudoku = document.querySelector('#sudoku');

    // Create board in table element.
    for (let row = 0; row < 9; row++) {
        let tr = document.createElement('tr');
        for (let col = 0; col < 9; col++) {
            let td = document.createElement('td');
            td.id = String.fromCharCode('a'.charCodeAt(0)+row)+(col+1);
            tr.appendChild(td);
        }
        sudoku.appendChild(tr);
    }

    // Create event handler for all cells.
    document.querySelectorAll('#sudoku tr td').forEach((td) => {
        td.addEventListener("mouseup", function(e) {
            let alreadyActive = td.classList.contains('active');
            document.querySelectorAll('#sudoku tr td.active').forEach((active) => {
                active.classList.remove('active');
            });
            if (!alreadyActive) this.classList.add('active');
        });
    });
}, false);
