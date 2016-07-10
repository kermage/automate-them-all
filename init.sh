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

cp $REAL_DIRECTORY/settings settings

read -p "Server Host: " SV_HOST
read -p "Server User: " SV_USER
read -p "Server Port: " SV_PORT
read -p "Database User: " DB_USER
read -p "Database Pass: " DB_PASS
read -p "Database Name: " DB_NAME

sed -i '/^SV_HOST/s/=.*/='$SV_HOST'/' settings
sed -i '/^SV_USER/s/=.*/='$SV_USER'/' settings
sed -i '/^SV_PORT/s/=.*/='$SV_PORT'/' settings
sed -i '/^DB_USER/s/=.*/='$DB_USER'/' settings
sed -i '/^DB_PASS/s/=.*/='$DB_PASS'/' settings
sed -i '/^DB_NAME/s/=.*/='$DB_NAME'/' settings