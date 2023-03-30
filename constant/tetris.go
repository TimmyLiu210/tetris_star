package constant

// initialize the tetris rotate rule
func Initialize(tetrisRotate *[TETRISTYPELENGTH][TETRIS_ROTATE_TYPE][TETRIS_X][TETRIS_Y]int){

	for tetrisType := range tetrisRotate{
		switch tetrisType{
		case TETRIS_I:
			tetrisRotate[TETRIS_I] = [TETRIS_ROTATE_TYPE][TETRIS_X][TETRIS_Y]int{
				//向右旋轉 直的最下面為 index [0]
				{
					{0,1,2,3},{0,-1,-2,-3},
				},
				{
					{0,-1,-2,-3},{0,1,2,3},
				},
			}
		case TETRIS_J:
			tetrisRotate[TETRIS_J] = [TETRIS_ROTATE_TYPE][TETRIS_X][TETRIS_Y]int{
				//向右旋轉 J的最左邊為 index [0]
				{
					{0,-1,0,1},{1,0,-1,-2},
				},
				{
					{1,0,-1,-2},{1,2,1,0},
				},
				{
					{1,2,1,0},{-2,-1,0,1},
				},
				{
					{-2,-1,0,1},{0,-1,0,1},
				},
			}		
		case TETRIS_L:
			tetrisRotate[TETRIS_J] = [TETRIS_ROTATE_TYPE][TETRIS_X][TETRIS_Y]int{
				//向右旋轉 L的最右邊為 index [0]
				{
					{-1,0,1,2},{0,1,0,-1},
				},
				{
					{0,1,0,-1},{2,1,0,-1},
				},
				{
					{2,1,0,-1},{-1,-2,-1,0},
				},
				{
					{-1,-2,-1,0},{-1,0,1,2},
				},
			}	
		case TETRIS_O:
		case TETRIS_S:
			tetrisRotate[TETRIS_S] = [TETRIS_ROTATE_TYPE][TETRIS_X][TETRIS_Y]int{
				//向右旋轉 S的最左下為 index [0]
				{
					{1,0,-1,-2},{0,1,0,1},
				},
				{
					{-1,0,1,2},{0,-1,0,-1},
				},
			}	
		case TETRIS_T:
			tetrisRotate[TETRIS_T] = [TETRIS_ROTATE_TYPE][TETRIS_X][TETRIS_Y]int{
				//向右旋轉 直的最下面為 index [0]
				{
					{1,0,-1,-1},{1,0,-1,1},
				},
				{
					{1,0,-1,1},{-2,-1,0,0},
				},
				{
					{-2,-1,0,0},{0,1,2,0},
				},
				{
					{0,1,2,0},{1,0,-1,-1},
				},
			}	
		case TETRIS_Z:
			tetrisRotate[TETRIS_Z] = [TETRIS_ROTATE_TYPE][TETRIS_X][TETRIS_Y]int{
				//向右旋轉 直的最下面為 index [0]
				{
					{0,-1,0,-1},{-1,0,1,2},
				},
				{
					{0,1,0,1},{1,0,-1,-2},
				},
			}	
		}
	}
	return
}