#!/bin/bash

#      ,#####    ###    ,#############
#     ###  `##  ## ##   ###      ##
#    ,##       ##   ##  ######   ##
#    ### ##### ######## #####'  ,##
#    ###   ######   `#####      ###
#    `#########'      ####     ####
#
# | kermage | PrivaTech -- GAFT | iMUT |
# Copyright 2016 Gene Alyson F. Torcende
# Email: genealyson.torcende@gmail.com
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# All rights reserved.


###    DO NOT CHANGE anything below    ###
### unless you know what you are doing ###
SCRIPT_SOURCE=${BASH_SOURCE[0]}
SCRIPT_SOURCE=${SCRIPT_SOURCE:-$0}
REAL_DIRECTORY=$( cd "$( dirname "$SCRIPT_SOURCE" )" && pwd )
LS=$( /bin/which ls )


function _gatas() {
	opts=$( $LS ${REAL_DIRECTORY} --ignore='gatas.*' )
	opts="${opts} init check"
	COMPREPLY=( $( compgen -W "${opts}" -- ${COMP_WORDS[COMP_CWORD]} ) )
}

function isFunction() {
	declare -F $1 &> /dev/null
	return $?
}

function init() {
	read -p "Server Host: " SV_HOST
	read -p "Server User: " SV_USER
	read -e -p "Server Port: " -i "22" SV_PORT
	read -e -p "Remote Path: " -i "~" RM_PATH
	read -p "Database Name: " DB_NAME
	read -p "Database User: " DB_USER
	read -p "Database Pass: " DB_PASS
	read -e -p "Database Host: " -i "localhost" DB_HOST
	read -e -p "Database Port: " -i "3306" DB_PORT

	SV_PORT="${SV_PORT:-22}"
	RM_PATH="${RM_PATH:-~}"
	DB_HOST="${DB_HOST:-localhost}"
	DB_PORT="${DB_PORT:-3306}"

	cp $REAL_DIRECTORY/gatas.cfg gatas.cfg
	sed -i '/^SV_HOST/s/=.*/='$SV_HOST'/' gatas.cfg
	sed -i '/^SV_USER/s/=.*/='$SV_USER'/' gatas.cfg
	sed -i '/^SV_PORT/s/=.*/='$SV_PORT'/' gatas.cfg
	sed -i '/^RM_PATH/s/=.*/='${RM_PATH//\//\\/}'/' gatas.cfg
	sed -i '/^DB_NAME/s/=.*/='$DB_NAME'/' gatas.cfg
	sed -i '/^DB_USER/s/=.*/='$DB_USER'/' gatas.cfg
	sed -i '/^DB_PASS/s/=.*/='$DB_PASS'/' gatas.cfg
	sed -i '/^DB_HOST/s/=.*/='$DB_HOST'/' gatas.cfg
	sed -i '/^DB_PORT/s/=.*/='$DB_PORT'/' gatas.cfg
}

function check() {
	SV_HOST=`grep '^SV_HOST' gatas.cfg | sed 's/[^=]*=//'`
	SV_USER=`grep '^SV_USER' gatas.cfg | sed 's/[^=]*=//'`
	SV_PORT=`grep '^SV_PORT' gatas.cfg | sed 's/[^=]*=//'`
	RM_PATH=`grep '^RM_PATH' gatas.cfg | sed 's/[^=]*=//'`
	DB_NAME=`grep '^DB_NAME' gatas.cfg | sed 's/[^=]*=//'`
	DB_USER=`grep '^DB_USER' gatas.cfg | sed 's/[^=]*=//'`
	DB_PASS=`grep '^DB_PASS' gatas.cfg | sed 's/[^=]*=//'`
	DB_HOST=`grep '^DB_HOST' gatas.cfg | sed 's/[^=]*=//'`
	DB_PORT=`grep '^DB_PORT' gatas.cfg | sed 's/[^=]*=//'`

	echo -n "Testing Local Database . . . "
	mysql --user="$DB_USER" --password="$DB_PASS" $DB_NAME -e "exit" > /dev/null 2>&1
	if [ "$?" -eq 0 ]; then
		echo "PASSED!"
	else
		echo "FAILED!"
	fi

	echo -n "Testing SSH Connection . . . "
	ssh $SV_USER@$SV_HOST -p$SV_PORT "exit" > /dev/null 2>&1
	if [ "$?" -eq 0 ]; then
		echo "PASSED!"
	else
		echo "FAILED!"
	fi

	echo -n "Testing Remote Path . . . "
	ssh $SV_USER@$SV_HOST -p$SV_PORT "test -w $RM_PATH" > /dev/null 2>&1
	if [ "$?" -eq 0 ]; then
		echo "PASSED!"
	else
		echo "FAILED!"
	fi

	echo -n "Testing Remote Database . . . "
	if [ "$DB_HOST" = "localhost" ]; then
		ssh $SV_USER@$SV_HOST -p$SV_PORT "mysql --user=\"$DB_USER\" --password=\"$DB_PASS\" $DB_NAME -e 'exit'" > /dev/null 2>&1
	else
		mysql --user="$DB_USER" --password="$DB_PASS" --host="$DB_HOST" --port="$DB_PORT" $DB_NAME -e "exit" > /dev/null 2>&1
	fi
	if [ "$?" -eq 0 ]; then
		echo 'PASSED!'
	else
		echo 'FAILED!'
	fi
}

function gatas() {
	COMMAND=$1
	CURRENT_PATH=$( pwd )
	GATAS_PATH=$CURRENT_PATH

	while [ "$GATAS_PATH" != "" ] && [ ! -e "$GATAS_PATH/gatas.cfg" ]; do
		GATAS_PATH=${GATAS_PATH%/*}
	done

	if [ "$GATAS_PATH" != "" ]; then
		cd $GATAS_PATH
	elif [ "$COMMAND" != "init" ]; then
		echo "Not a GATAS-ready directory (or any of the parent directories)"
	fi

	if [ ! -r "gatas.cfg" ] && [ "$COMMAND" != "init" ]; then
		echo "Configuration file 'gatas.cfg' does not exist. Run 'gatas init' first."
		return 1
	fi

	if ( isFunction "$COMMAND" ); then
		"$COMMAND"
	else
		COMMAND_FILE=$( $LS ${REAL_DIRECTORY}/${COMMAND} 2> /dev/null | head -1 )

		if [ -r "$COMMAND_FILE" ]; then
			bash "$COMMAND_FILE" "${@: 2}"
		else
			echo "'$COMMAND' is not a gatas command. Use 'gatas <command>':"
			$LS ${REAL_DIRECTORY} --ignore='gatas.*'
			return 1
		fi
	fi

	cd $CURRENT_PATH
}

complete -F _gatas gatas
