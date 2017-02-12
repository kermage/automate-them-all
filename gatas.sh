#!/system/bin/sh

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
REAL_DIRECTORY=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

function isFunction() {
    declare -F $1 &> /dev/null
    return $?
}

function init() {
    cp $REAL_DIRECTORY/gatas.cfg gatas.cfg
    
    read -p "Server Host: " SV_HOST
    read -p "Server User: " SV_USER
    read -p "Server Port: " SV_PORT
    read -p "Remote Path: " RM_PATH
    read -p "Database User: " DB_USER
    read -p "Database Pass: " DB_PASS
    read -p "Database Name: " DB_NAME
    
    sed -i '/^SV_HOST/s/=.*/='$SV_HOST'/' gatas.cfg
    sed -i '/^SV_USER/s/=.*/='$SV_USER'/' gatas.cfg
    sed -i '/^SV_PORT/s/=.*/='$SV_PORT'/' gatas.cfg
    sed -i '/^RM_PATH/s/=.*/='${RM_PATH//\//\\/}'/' gatas.cfg
    sed -i '/^DB_USER/s/=.*/='$DB_USER'/' gatas.cfg
    sed -i '/^DB_PASS/s/=.*/='$DB_PASS'/' gatas.cfg
    sed -i '/^DB_NAME/s/=.*/='$DB_NAME'/' gatas.cfg
}

function test() {
    SV_HOST=`grep '^SV_HOST' gatas.cfg | sed 's/.*=//'`
    SV_USER=`grep '^SV_USER' gatas.cfg | sed 's/.*=//'`
    SV_PORT=`grep '^SV_PORT' gatas.cfg | sed 's/.*=//'`
    DB_USER=`grep '^DB_USER' gatas.cfg | sed 's/.*=//'`
    DB_PASS=`grep '^DB_PASS' gatas.cfg | sed 's/.*=//'`
    DB_NAME=`grep '^DB_NAME' gatas.cfg | sed 's/.*=//'`
    
    echo -n "Testing Local Database . . . "
    mysql --user=$DB_USER --password=$DB_PASS $DB_NAME -e "exit" > /dev/null 2>&1
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
    
    echo -n "Testing Remote Database . . . "
    ssh $SV_USER@$SV_HOST -p$SV_PORT "mysql --user=$DB_USER --password=$DB_PASS $DB_NAME -e 'exit'" > /dev/null 2>&1
    if [ "$?" -eq 0 ]; then
        echo 'PASSED!'
    else
        echo 'FAILED!'
    fi
}

if [ $# -gt 0 ]; then
    COMMAND=$1
    
    if ( isFunction "$COMMAND" ); then
        "$COMMAND"
    else
        COMMAND_FILE=$( ls ${REAL_DIRECTORY}/${COMMAND} 2> /dev/null | head -1 )
        
        if [ -r "$COMMAND_FILE" ]; then
            sh "$COMMAND_FILE" "${@: 2}"
        fi
    fi
fi
