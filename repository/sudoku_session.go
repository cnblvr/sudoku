package repository

//type SudokuSession struct {
//	conn     redis.Conn
//	id       uuid.UUID
//	sudokuID int64
//}
//
//func SudokuSessionByID(conn redis.Conn, id uuid.UUID) (SudokuSession, error) {
//	sudokuID, err := redis.Int64(conn.Do("GET", keySudokuSession(id)))
//	switch err {
//	case nil:
//	case redis.ErrNil:
//		return SudokuSession{}, nil
//	default:
//		return SudokuSession{}, err
//	}
//	return SudokuSession{
//		conn:     conn,
//		id:       id,
//		sudokuID: sudokuID,
//	}, nil
//}
//
//func SudokuSessionByIDString(conn redis.Conn, idStr string) (SudokuSession, error) {
//	id, err := uuid.FromString(idStr)
//	if err != nil {
//		return SudokuSession{}, err
//	}
//	return SudokuSessionByID(conn, id)
//}
//
//func NewSudokuSession(conn redis.Conn, sudoku Sudoku, user User) (SudokuSession, error) {
//	if sudoku.IsNull() {
//		return SudokuSession{}, fmt.Errorf("sudoku is null")
//	}
//	sudokuSessionSeed, err := redis.Int64(conn.Do("INCR", keyLastSudokuSessionSeed()))
//	if err != nil {
//		return SudokuSession{}, err
//	}
//	hash := md5.Sum([]byte("sudoku_session" + strconv.FormatInt(sudokuSessionSeed, 16)))
//	id, err := uuid.FromBytes(hash[:])
//	if err != nil {
//		return SudokuSession{}, err
//	}
//
//	if _, err := conn.Do("SET", keySudokuSession(id), sudoku.id); err != nil {
//		return SudokuSession{}, err
//	}
//	if !user.IsNull() {
//		if _, err := conn.Do("SET", keySudokuSessionUserID(id), user.id); err != nil {
//			return SudokuSession{}, err
//		}
//	}
//
//	return SudokuSession{
//		conn: conn,
//		id:   id,
//	}, nil
//}
//
//func (s SudokuSession) Sudoku() Sudoku {
//	return Sudoku{
//		conn: s.conn,
//		id:   s.sudokuID,
//	}
//}
//
//func (s SudokuSession) ID() uuid.UUID {
//	return s.id
//}
//
//func (s SudokuSession) IsNull() bool {
//	return s.id.String() == "00000000-0000-0000-0000-000000000000"
//}
//
//func keySudokuSession(id uuid.UUID) string {
//	return fmt.Sprintf("sudoku_session:%s", id.String())
//}
//
//func keyLastSudokuSessionSeed() string {
//	return fmt.Sprintf("last_sudoku_session_seed")
//}
//
//func keySudokuSessionUserID(id uuid.UUID) string {
//	return fmt.Sprintf("%s:user_id", keySudokuSession(id))
//}
