#!/usr/bin/env bash
set -e
set -x

DINGDING_ROBOT_ADDR="https://oapi.dingtalk.com/robot/send?access_token=4dd14890111f3a23928531b1e1f290f1669ec227b813e68c3c7ee09dfb4000a9"

UNZIP_COMMAND="unzip -o " 
ZIP_COMMAND="zip -q -r -l -D "

TEMP_ZIP_NAME="Temp.zip"
TEMP_NAME="Temp"

#国外资源同步指令
UNITY3D_IN_SYNC_CMD="aws s3 cp "

#国外的CDN路径
UNITY3D_OUT_PATH="s3://diablo-dev"

if [ ! -f ${TEMP_ZIP_NAME} ]; then
    exit 0
fi

${UNZIP_COMMAND} ${TEMP_ZIP_NAME} -d ${TEMP_NAME}


if [ "$?"-ne 0]; then
    exit 0 
fi


if [ ! -d ${TEMP_NAME} ]; then
    exit 0
fi

rm -f ${TEMP_ZIP_NAME}

cv=0
rv=0
rv_init=0
channel="GP"
target="Android"
for cvs in `cat ${TEMP_NAME}/Version.txt`
do
    channel="$(cut -d'.' -f1 <<<"${cvs}")"
    cv="$(cut -d'.' -f2 <<<"${cvs}")"
    rv="$(cut -d'.' -f3 <<<"${cvs}")"
    rv_init="$(cut -d'.' -f4 <<<"${cvs}")" 
    target="$(cut -d'.' -f5 <<<"${cvs}")" 
done

echo ${channel}.${rv}.${rv_init}.${target}

if [[ ${rv} == 0 ]]; then
    curl ${DINGDING_ROBOT_ADDR} \
   -H 'Content-Type: application/json' \
   -d "{\"msgtype\": \"markdown\", 
       \"markdown\": {
          \"title\":\"[开发]新更新事件\",
            \"text\": \"###无法或者版本号 \n \"
        }
      }"
    exit 0
fi


if [ ! -d "${channel}" ]; then 
    mkdir ${channel}
fi
if [ ! -d "${channel}/${target}" ]; then 
    mkdir "${channel}/${target}"
fi
if [ ! -d "${channel}/${target}/${rv}" ]; then 
    mkdir "${channel}/${target}/${rv}"
fi
cd ${TEMP_NAME}
for (( i = ${rv_init}; i < ${rv}; i++ )); do

if [ ! -f "../${channel}/${target}/${rv}/update_${i}.zip" ]
  then    
    #statements
    for (( j = ${i} + 1; j <= ${rv}; j++ )); do
      #statements
      if [ ! -d "../${channel}/${target}/${rv}/${i}" ]
      then
        mkdir ../${channel}/${target}/${rv}/${i}  
      fi  
      
      if [ -d ${j} ]
      then
        cp -r ${j} ../${channel}/${target}/${rv}/${i}
      fi
    done

    cd ../${channel}/${target}/${rv}/${i}

    ${ZIP_COMMAND} ../update_${i}.zip .

    size="$(du -k ../update_${i}.zip)"

    echo -n ${size}>../info_${i}.txt
    cd ..
    rm -rf ${i}
    
    ${UNITY3D_IN_SYNC_CMD} update_${i}.zip ${UNITY3D_OUT_PATH}/${channel}/${target}/${rv}/update_${i}.zip
    ${UNITY3D_IN_SYNC_CMD} info_${i}.txt ${UNITY3D_OUT_PATH}/${channel}/${target}/${rv}/info_${i}.txt
    rm -f update_${i}.zip
    cd ../../../${TEMP_NAME}
  else
    ${UNITY3D_IN_SYNC_CMD} ../${channel}/${target}/${rv}/update_${i}.zip ${UNITY3D_OUT_PATH}/${channel}/${target}/${rv}/update_${i}.zip
    ${UNITY3D_IN_SYNC_CMD} ../${channel}/${target}/${rv}/info_${i}.txt ${UNITY3D_OUT_PATH}/${channel}/${target}/${rv}/info_${i}.txt
    rm -f ../${channel}/${target}/${rv}/update_${i}.zip
  fi
  
 
 if [ "$?"-ne 0]; then 
    curl ${DINGDING_ROBOT_ADDR} \
   -H 'Content-Type: application/json' \
   -d "{\"msgtype\": \"markdown\", 
       \"markdown\": {
          \"title\":\"[开发]新更新事件\",
            \"text\": \"###[${channel}][${cv}][${rv}][${i}]上传失败 \n \"
        }
      }"
 exit 1; 
 fi

done
cd ..
rm -rf ${TEMP_NAME}

curl ${DINGDING_ROBOT_ADDR} \
   -H 'Content-Type: application/json' \
   -d "{\"msgtype\": \"markdown\", 
       \"markdown\": {
          \"title\":\"[开发]新更新事件\",
            \"text\": \"###[${channel}][${target}]热更版本 ${cv}.${rv}:创建成功 \n \"
        }
      }"
    exit 0