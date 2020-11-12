FSN - File System Notifier
==========================

Описание
--------
FSN предназначен для выполнения команд при возникновении событий в файловой системе. Может наблюдать как за отдельными файлами, так и за каталогами.

Функции
-------
1. Наблюдает за файлами и каталогами.
2. Реагирует на события: create, write, remove, rename, chmod.
3. Пауза перед выполнением команды.
4. Установка переменных среды для команды.
5. Установка рабочего каталога для команды.
6. Подавление последующих событий в течение указанного времени.

Параметры командной строки
---------------------------
    fsn -config <path to config file>

- **config** -- путь к файлу конфигурации. По умолчанию: **fsn.cfg** в текущем каталоге.

Файл конфигурации
-----------------

Файл конфигурации записывается в формате [JSON](https://ru.wikipedia.org/wiki/JSON).

    {
        "workers": 2,
        "watch": {
            "/home/serg/work/fsn" : {
                "cmd": "/usr/bin/bash echo 111",
                "pause": 0,
                "modes": ["create", "write"],
                "log_output": true,
                "throttle": 1500,
                "workdir": "/tmp",
                "env": ["PATH=/tmp"]
            }
        }
    }

- **workers** - количество процессов, выполняющих команды. Количество процессов большее 1 может ускорить выполнение команд, так как в случае поступления большого количества событий команды смогут выполняться параллельно.
- **watch** - наблюдаемые объекты файловой системы. В качестве ключа указывается путь, в качестве значения -- объект с параметрами выполнения команды. Пути могут быть как непосредственно файлами, так и каталогами. **В случае наблюдения за каталогом события регистрируются для всех файлов внутри каталога, но не в подкаталогах.**
- **cmd** - выполняемая команда. В ОС Windows при необходимости выполнить командный файл (.bat, .cmd) необходимо указать командный интерпретатор, например cmd.exe. Пример:

        "cmd.exe /q /c c:\\tools\\test.bat"

- **pause** - пауза в миллисекундах перед выполнением команды.
- **modes** - список, содержащий операции при возникновении которых выполняется команда. При пустом списке ([]), команда выполняется при наступлении любого события. Возможные варианты: "create", "write", "remove", "rename", "chmod".
- **log_output** - флаг, указывающий на необходимость вывода результата выполнения команды в лог. Выводятся stdout и stderr.
- **throttle** - время в миллисекундах в течение которого игнорируются последующие события для данного пути.
- **workdir** - каталог в котором будет выполняться программа.
- **env** - список переменных среды. Каждая переменная описывается одной строкой в виде "переменная=значение". В начало списка всегда добавляются переменные текущего процесса. Переменные могут повторяться, в этом случае используется послднее указанное значение.

В случае работы под Windows слеши в путях необходимо экранировать. Пример:
    
    "C:\\tools\\bats\\x.bat"

Запустить программу в качестве сервиса Windows можно с помощью [NSSM](https://nssm.cc/).

## Licenses

This software uses following libraries, which have its own licenses.

### [fsnotify](https://github.com/fsnotify/fsnotify)

Copyright (c) 2012 The Go Authors. All rights reserved.

Copyright (c) 2012-2019 fsnotify Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

 ### [go-shlex](https://github.com/anmitsu/go-shlex)

 Copyright (c) anmitsu <anmitsu.s@gmail.com>