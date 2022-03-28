'use strict';

class Sudoku {
    #_object;
    #_keyboard;
    #isWin = false;
    #gameID;
    #ws;

    constructor(param) {
        if (!param)
            throw 'sudoku: parameters not defined';
        if (!param.selector || typeof param.selector !== 'string')
            throw 'sudoku: required parameter \'selector\' is not defined or not string';
        if (!param.gameID || typeof param.gameID !== 'string')
            throw 'sudoku: required parameter \'gameID\' is not defined or not string';
        this.#gameID = param.gameID;
        this.#_object = document.querySelector(param.selector);
        if (!this.#_object)
            throw 'sudoku: object by parameter \'selector\' not found';
        if (param.allowEditing) {
            if (typeof param.allowEditing !== 'boolean')
                throw 'sudoku: parameter \'allowEditing\' is not boolean';
            if (param.keyboardSelector) {
                if (typeof param.keyboardSelector !== 'string')
                    throw 'sudoku: parameter \'keyboardSelector\' is not string';
                this.#_keyboard = document.querySelector(param.keyboardSelector);
                if (!this.#_keyboard)
                    throw 'sudoku: object by parameter \'keyboardSelector\' not found';
            }
        }

        for (let row = 0; row < 9; row++) {
            let _row = document.createElement('div');
            _row.classList.add('sud-row');
            for (let col = 0; col < 9; col++) {
                let _cell = document.createElement('div');
                _cell.classList.add('sud-cll', 'is-cnd');
                // create digit field
                let _dgt = document.createElement('div');
                _dgt.classList.add('sud-dgt');
                _cell.appendChild(_dgt);
                // create table of candidates
                let _cnd = document.createElement('div');
                _cnd.classList.add('sud-cnd');
                for (let idx = 1; idx <= 9; idx++) {
                    let _cndItem = document.createElement('span');
                    _cndItem.textContent = ''+idx;
                    _cnd.appendChild(_cndItem);
                }
                _cell.appendChild(_cnd);
                _row.appendChild(_cell);
            }
            this.#_object.appendChild(_row);
        }

        if (this.#_keyboard) {
            let createBtn = (label, event) => {
                let btn = document.createElement('div');
                btn.classList.add('kb-btn');
                btn.textContent = label;
                btn.addEventListener('click', event);
                this.#_keyboard.appendChild(btn);
            }
            createBtn( 'c', (e) => {
                console.log('TODO set candidates');
            });
            createBtn( '⨯', (e) => {
                this.#placeDigitInActive('0');
            });
            for (let digit = 1; digit <= 9; digit++) {
                createBtn(digit, (e) => {
                    this.#placeDigitInActive(''+digit);
                });
            }
        }

        if (param.allowEditing) {
            this.#_object.querySelectorAll('.sud-cll').forEach((_cell) => {
                _cell.addEventListener('mouseup', (e) => {
                    this.#setActive(_cell);
                });
            });

            document.addEventListener('keyup', (e) => {
                if (e.defaultPrevented) {
                    return;
                }
                let _a = this.#_object.querySelector('.sud-cll.active');
                switch (e.code) {
                    case 'ArrowUp':    this.#setActive(_a, 'up');    break;
                    case 'ArrowRight': this.#setActive(_a, 'right'); break;
                    case 'ArrowDown':  this.#setActive(_a, 'down');  break;
                    case 'ArrowLeft':  this.#setActive(_a, 'left');  break;
                    case 'Digit0':
                    case 'Numpad0':
                    case 'Space':
                    case 'Backspace':
                        this.#placeDigitInActive('0');
                }
                if ('1' <= e.key && e.key <= '9') {
                    this.#placeDigitInActive(e.key);
                }
            });
        }

        this.#_object.addEventListener('apiReady', () => {
            this.#ws.send('getPuzzle', {
                game_id: this.#gameID,
                need_candidates: true,
            });
        }, {once: true});

        this.#_object.addEventListener('api_getPuzzle', (e) => {
            let body = e.detail.body;
            this.#_object.querySelectorAll('.sud-row').forEach((_row, row) => {
                _row.querySelectorAll('.sud-cll').forEach((_cell, col) => {
                    this.#placeDigit(_cell, '0', true);
                    let d = body.puzzle[row*9+col];
                    if ('1' <= d && d <= '9') {
                        this.#placeDigit(_cell, d, true);
                        _cell.classList.add('hint');
                    }
                    if (body.candidates) {
                        this.#setCandidatesFor(_cell, body.candidates[this.#stringifyPoint(row, col)]);
                    }
                });
            });
        });

        this.#_object.addEventListener('api_makeStep', (e) => {
            let body = e.detail.body;
            this.#_object.querySelectorAll('.sud-cll').forEach((_cell) => {
                _cell.classList.remove('error');
            });
            if (body.win) {
                this.#isWin = true;
                alert('win'); // TODO
                return;
            }
            body.errors = this.#parsePoints(body.errors);
            this.#_object.querySelectorAll('.sud-row').forEach((_row, row) => {
                _row.querySelectorAll('.sud-cll').forEach((_cell, col) => {
                    body.errors.forEach((p) => {
                        if (p.row === row && p.col === col) {
                            _cell.classList.add('error');
                        }
                    });
                    if (body.candidates) {
                        this.#setCandidatesFor(_cell, body.candidates[this.#stringifyPoint(row, col)]);
                    }
                });
            });
        });
    }

    dispatchEvent(ce) {
        this.#_object.dispatchEvent(ce);
    }

    connectWS(ws) {
        this.#ws = ws;
    }

    #placeDigit(_cell, digit, notMakeStep) {
        if (this.#isWin) return;
        if (!_cell || _cell.classList.contains('hint')) return;
        let _digit = _cell.querySelector('.sud-dgt');
        let oldDigit = _digit.textContent===''?'0':_digit.textContent;
        _cell.classList.remove('is-dgt', 'is-cnd');
        if (digit === '0') {
            _digit.textContent = '';
            _cell.classList.add('is-cnd');
        } else {
            _digit.textContent = digit;
            _cell.classList.add('is-dgt');
        }
        if (!notMakeStep && oldDigit !== digit) this.#apiMakeStep();
    }

    #placeDigitInActive(digit, notMakeStep) {
        if (this.#isWin) return;
        this.#placeDigit(this.#_object.querySelector('.sud-cll.active'), digit, notMakeStep);
    }

    #setCandidatesFor(_cell, cands) {
        if (!_cell || !cands) return;
        _cell.querySelectorAll('.sud-cnd span').forEach((_span) => {
            if (cands.includes(_span.textContent.charCodeAt(0)-'0'.charCodeAt(0))) {
                _span.classList.remove('hidden');
            } else {
                _span.classList.add('hidden');
            }
        });
    }

    #setActive(_cell, dir) {
        if (!_cell) {
            _cell = this.#_object.querySelectorAll('.sud-row').item(9/2).querySelectorAll('.sud-cll').item(9/2);
            dir = undefined;
            if (!_cell) return;
        }
        let _row = _cell.closest('.sud-row');
        switch (dir) {
            case 'up':
                let _prev = _row.previousElementSibling;
                if (!_prev) return;
                _cell = _prev.querySelectorAll('.sud-cll').item(this.#getIndex(_cell));
                break;
            case 'right':
                _cell = _cell.nextElementSibling; break;
            case 'down':
                let _next = _row.nextElementSibling;
                if (!_next) return;
                _cell = _next.querySelectorAll('.sud-cll').item(this.#getIndex(_cell));
                break;
            case 'left':
                _cell = _cell.previousElementSibling; break;
        }
        if (!_cell) return;
        let isAlready = _cell.classList.contains('active');
        this.#_object.querySelectorAll('.sud-cll.active').forEach((_active) => {
            _active.classList.remove('active');
        });
        if (!isAlready) _cell.classList.add('active');
    }

    #apiMakeStep() {
        let state = '';
        this.#_object.querySelectorAll('.sud-dgt').forEach((_dgt) => {
            let val = _dgt.textContent;
            if (val === '') val = '.';
            state += val;
        });
        this.#ws.send('makeStep', {
            game_id: this.#gameID,
            state: state,
            need_candidates: true,
        })
    }

    #getIndex(_node) {
        let index = 0;
        while (_node = _node.previousElementSibling) {
            index++;
        }
        return index;
    }

    #parsePoints(points) {
        let out = [];
        if (!points) return out;
        points.forEach((p) => {
            out = out.concat([{
                row: p[0].charCodeAt(0)-'a'.charCodeAt(0),
                col: parseInt(p[1])-1,
            }]);
        });
        return out;
    }

    #stringifyPoint(row, col) {
        return String.fromCharCode((row)+'a'.charCodeAt(0)) + (col+1);
    }
}
