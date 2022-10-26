# service_xml_json
main.go:
Парсинг xml файла.
На вход принимает base64 - закодированный файл xml. Отдает строку с JSON.

Ход установки на ВМ:

Создаем директорию:
`mkdir $GOPATH/parsexml && cd parsexml`
`nano main.go`
Помещаем код из main.go
`ctrl+o ctrl+x`

Далее создаем модуль (example.com/m - если не нужен репо)
`go mod init example.com/m`
Устанавливаем все необходимое
`go mod tidy`
Создаем билд
`go build main.go`

Далее идем создавать `systemd unit file`, который фоном запустит наш сервер:
Переходим в директорию: `cd /lib/systemd/system`
Создаем файл: `nano parsexml.service`
Прописываем в этом файле следующее:
```
[Unit]
Description=parsexml

[Service]
Type=simple
Restart=always
RestartSec=5s
ExecStart=/parsexml/main

[Install]
WantedBy=multi-user.target
```
Сохраняем файл (ctrl + O) и Выходим (ctrl + X)
Далее запускаем сервис: `service parsexml start`
Проверяем статус (должен подсветиться зеленым и Active): `service parsexml status`
Теперь по адресу: `http://${АДРЕС ВМ}:6060/parse-xml`

Выполняем запрос на этот адрес:
```
// Получаем файл xml и преобразуем его в  строку вида base64
const xml_file = await Context.data.xml_file!.fetch()                        
const fileObj = await fetch(await xml_file!.getDownloadUrl());
const content = new Uint8Array(await fileObj.arrayBuffer());
let binary = '';
for (const char of content) {
     binary += String.fromCharCode(char);
}
const base64 = btoa(binary);

const res = await fetch(`http://${АДРЕС ВМ}:6060/parse-xml`, {
     method: "POST",
     body: base64
});

const answer: any = await res.json()
```
