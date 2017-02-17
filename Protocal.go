package chathome

const (
	/**
	 *  SPLIT SYMBOLS
	 */
	P_SP		= ":::"
	P_SP_SEND	= "->"
	P_SP_RCV	= "<-"
	P_SP_ARG	= ","

	/**
	 *  LOW GRADE CALL
	 */
	P_SEND_MSG 	= "MSG"
	P_SEND_G_MSG	= "MSG_G"
	P_SEND_BR	= "BROADCAST"
	P_CALL		= "CALL"
	P_RESULT	= "RESULT"

	/**
	 *  METHODS
	 *  call is for server
	 *  result is for client
	 */

	P_CALL_LOGIN		= "LOGIN"
	P_CALL_REG		= "REGISTER"
	P_CALL_QUIT		= "QUIT"
	P_CALL_USER_LIST	= "USER_LIST"

	P_RS_LOGIN		= "LOGIN"
	P_RS_REG		= "REGISTER"
	P_RS_USER_LIST		= "USER_LIST"

	/**
	 *  Errors
	 */
	P_RS_SUCCESS	= "SUCCESS"
	P_RS_ERR	= "ERROR:"
	E_CODE_EXISTS	= "101"		//USER IS EXISTS
	E_CODE_PWD	= "102"		//PWD IS WRONG

)