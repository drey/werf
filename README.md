# dapp [![Gem Version](https://badge.fury.io/rb/dapp.svg)](https://badge.fury.io/rb/dapp) [![Build Status](https://travis-ci.org/flant/dapp.svg)](https://travis-ci.org/flant/dapp) [![Code Climate](https://codeclimate.com/github/flant/dapp/badges/gpa.svg)](https://codeclimate.com/github/flant/dapp) [![Test Coverage](https://codeclimate.com/github/flant/dapp/badges/coverage.svg)](https://codeclimate.com/github/flant/dapp/coverage)

## Reference

### Основные определения

#### Проект
Проект (project) — это директория, содержащая приложение или набор приложений (см. [приложение](#Приложение)).
* Приложение может находиться в корне проекта (в этом случае в корне проекта лежит соответствующий Dappfile).
* В случае, если в проекте есть несколько приложений — они находятся в директориях .dapps/\<имя-приложения\>/ (в каждой из которых есть соответствующий Dappfile).

#### Директория проекта
Директория проекта (project dir) — это директория, в которой находится директория .dapps или, при отсутствии .dapps — это директория, содержащая Dappfile.

#### Имя проекта
Имя проекта — это последний элемент пути к git репозиторию из параметра конфигурации remote.origin.url или, при отсутствии git или параметра конфигурации remote.origin.url — имя директории корня проекта.

#### Dappfile
Dappfile — это файл, содержащий инструкции по сборке docker образов приложения (см. [приложение](#Приложение)).

#### Приложение
Приложение (dapp) — это набор правил, объединенных в одном Dappfile, по которым происходит сборка одного или нескольких подприложений.
* В рамках одного приложения может быть описано дерево подприложений.
* При сборке дерева подприложений, docker образы будут собраны для всех подприложений листьев описанного дерева.

#### Директория приложения
Директория приложения (home dir) — это директория, содержащая [Dappfile](#Dappfile) [приложения](#Приложение).

#### Базовое имя приложения
Базовое имя приложения (basename) ­— это имя, связанное с каждым Dappfile.
* По умолчанию базовое имя приложения ­— это имя директории, в которой находится Dappfile.
* Базовое имя приложения может быть переопределено в Dappfile (см. [name](#name-name)).

#### Подприложение
Подприложение (app) — это средство группировки правил сборки в иерархию с наследованием.
* Подприложение наследует правила сборки того подприложения, в котором оно объявлено, и глобальные правила сборки.
* Сборка docker образов осуществляется только для тех подприложений, которые являются листьями в описанном дереве.
* Для каждого подприложения в Dappfile указывается имя.
  * При этом вложенные подприложения наследуют имена родительских подприложений и базовое имя приложения (см. [базовое имя приложения](#Базовое-имя-приложения)).
  * Итоговое имя подприложения имеет вид: \<базовое имя приложения\>-\<подприложение-1\>-\<подприложение-2\>...-\<подприложение-N\>.
  
#### Артефакт
Артефакт (artifact) — это несамостоятельное [приложение](#Приложение).
* Используется для изолирования процесса сборки (среды, программного обеспечение, данных) необходимых ресурсов от связанных приложений.
* Наследования от родителя происходит по тем же правилам, что и у [подприложения](#Подприложение).
* Стадии артефакта практически полностью совпадают со стадиями приложения, но:
  * Отсутствуют:
    * [git_artifact_post_setup_patch](#git_artifact_post_setup_patch).
    * [git_artifact_latest_patch](#git_artifact_latest_patch).
    * [docker_instructions](#docker_instructions).
  * Добавляются:
    * [git_artifact_artifact_patch](#git_artifact_artifact_patch).
    * [build_artifact](#build_artifact).
* Могут использоваться такие же директивы, как у приложения, но с некоторыми особенностями:
  * [docker](#docker): не используются директивы применения dockerfile-инструкций.
  * [shell](#shell): добавляется директива [build_artifact](#shellbuild_artifact--cache_version-cache_version).
  * У артефакта не может быть потомков, поэтому отсутствуют директивы app и artifact.

#### Стадия
Стадия (stage) — это именованный набор инструкций для сборки docker образа.
* Собранное приложение представляет собой цепочку связанных стадий.
* Имя docker образа стадии формируется по шаблону: dappstage-\<[имя проекта](#Имя-проекта)\>-\<[базовое имя приложения](#Базовое-имя-приложения)\>:\<[сигнатура стадии](#Сигнатура-стадии)\>.

#### Сигнатура стадии
Сигнатура стадии (signature) — это контрольная сумма правил сборки, зависимостей стадии и сигнатуры предыдущей стадии, если она существует.
* Изменение сигнатуры стадии ведет к пересборке этой стадии.
* При отсутствии правил и зависимостей, стадия игнорируется, используется сигнатура предыдущей стадии.

#### Кэш приложения
Набор docker-образов, который определяется [стадиями приложения](#Стадия).

#### Тип сборки
Dapp поддерживает 2 типа сборки: [chef приложение](#Chef-приложение) и [shell приложение](#Shell-приложение).
* В одном [приложении](#Приложение) может быть использован один тип сборки.
* Тип сборки определяется автоматически при первом использовании соответствующей инструкции в [Dappfile](#Dappfile).

#### Chef приложение
Chef приложение — это приложение, для сборки которого используются [cookbook приложения](#cookbook-приложения) и [mdapp модули](#mdapp-модуль).

#### Cookbook приложения
Cookbook приложения (cookbook dapp) — это основной chef cookbook, связанный с одним [приложением](#Приложение).
* Обязательные файлы:
  * metadata.rb
    * Содержит все зависимости, в т.ч. [mdapp модули](#mdapp-модуль).
  * Berksfile
    * Содержит все зависимости, в т.ч. [mdapp модули](#mdapp-модуль).
  * Berksfile.lock
    * Может отсутствовать в [dev режиме](#Режим-разработки) — будет создан во время выполнения сборки.
* Структура файлов.
  * Атрибуты.
    * Файлы атрибутов не поддерживаются.
    * Хэш normal-атрибутов задается через [Dappfile](#dappfile), см. [chef.attribute](#chefattribute), [chef.\<стадия\>\_attribute](#chefстадия_attribute).
  * Файлы.
    * Директория files/\<стадия\>/\<рецепт\> — содержит файлы рецепта для стадии, опциональна.
    * Директория files/\<стадия\>/common — содержит общие файлы для стадии, опциональна.
    * Недопустимо наличие одинаковых имен файлов в директории рецепта и common.
  * Шаблоны.
    * Директория templates/\<стадия\>/\<рецепт\> — содержит файлы шаблонов рецепта для стадии, опциональна.
    * Директория templates/\<стадия\>/common — содержит общие файлы шаблонов для стадии, опциональна.
    * Недопустимо наличие одинаковых имен файлов в директории рецепта и common.
  * Рецепты.
    * Директория recipes/\<стадия\> — содержит файлы рецептов для стадии, опциональна.

#### Mdapp модуль
Mdapp модуль — это дополнительный chef cookbook, который подключается к сборке [chef приложения](#chef-приложение).
* Обязательные файлы:
  * metadata.rb
* Структура файлов.
  * Файлы атрибутов не поддерживаются.
  * Файлы.
    * Директория files/\<стадия\> — содержит файлы доступные для стадии, опциональна.
    * Директория files/common — содержит файлы доступные для всех стадий, опциональна.
    * Недопустимо наличие одинаковых имен файлов в директории стадии и common.
  * Шаблоны.
    * Директория template/\<стадия\> — содержит файлы доступные для стадии, опциональна.
    * Директория template/common — содержит файлы доступные для всех стадий, опциональна.
    * Недопустимо наличие одинаковых имен файлов в директории стадии и common.
  * Рецепты.
    * Файл recipes/\<стадия\>.rb — файл рецепта модуля для стадии, опционален.

#### Стадия cookbook\`а
Стадия cookbook\`а — это часть cookbook\`а, которая используется при сборке стадии для cookbook\`а приложения и mdapp модулей.
* Понятие применимо только к cookbook\`у приложения и mdapp модулям.
* Для всех остальных cookbook\`ов при сборке стадии используется все файлы cookbook\`а.

#### Установка стадии cookbook\`а
Установка стадии cookbook\`а — это процесс копирования файлов стадии cookbook\`а во [временное хранилище](#Временная-директория-приложения), подключаемое в дальнейшем в контейнер для сборки стадии.
* Установка [cookbook\`а приложения](#cookbook-приложения).
  * Атрибуты.
    * Хэш атрибутов, совмещенных из [chef.attribute](#chefattribute) и [chef.\<стадия\>\_attribute](#chefстадия_attribute), устанавливается в normal-атрибуты, передаваемые через JSON-файл.
    * Файлы атрибутов не поддерживаются.
  * Файлы.
    * Содержимое директории files/\<стадия\>/common при наличии устанавливается в директории files/default.
    * Содержимое директории files/\<стадия\>/\<рецепт\> устанавливается в директорию files/default.
      * Для каждого включенного рецепта при наличии соответствующей директории.
    * Недопустимо наличие одинаковых имен файлов в директории рецепта и common.
  * Шаблоны.
    * Содержимое директории templates/\<стадия\>/common при наличии устанавливается в директории templates/default.
    * Содержимое директории templates/\<стадия\>/\<рецепт\> устанавливается в директорию templates/default.
      * Для каждого включенного рецепта при наличии соответствующей директории.
    * Недопустимо наличие одинаковых имен файлов в директории рецепта и common.
  * Рецепты.
    * Файл recipes/\<стадия\>/\<рецепт\>.rb устанавливается в recipes/\<рецепт\>.rb.
      * Для каждого включенного рецепта при его наличии.
    * При отсутствии рецептов генерируется пустой рецепт recipes/void.rb.
      * Отсутствие рецептов подразумевает одно из условий:
        * отсутствие включенных рецептов в конфигурации;
        * отсутствие файлов рецептов (recipes/\<стадия\>/\<рецепт\>.rb) для всех включенных в конфигурации рецептов.
      * Это позволяет активировать атрибуты, объявленные в данном cookbook\`е.
* Установка [mdapp модуля](#mdapp-модуль).
  * Файлы атрибутов не поддерживаются.
  * Файлы.
    * Содержимое директории files/common при наличии устанавливается в директорию files/default.
    * Содержимое директории files/\<стадия\> при наличии устанавливается в директорию files/default.
    * Недопустимо наличие одинаковых имен файлов в директории стадии и common.
  * Шаблоны.
    * Содержимое директории templates/common при наличии устанавливается в директорию templates/default.
    * Содержимое директории templates/\<стадия\> при наличии устанавливается в директорию templates/default.
    * Недопустимо наличие одинаковых имен файлов в директории стадии и common.
  * Рецепты.
    * Файл recipes/\<стадия\>.rb при наличии устанавливается в recipes/\<стадия\>.rb.
    * При отсутствии рецепта генерируется пустой рецепт recipes/void.rb.
      * Это позволяет активировать атрибуты, объявленные в данном cookbook\`е, при отсутствии рецепта.
* Остальныe cookbook\`и устанавливаются без изменений, "как есть".

#### Контрольная сумма cookbook\`ов стадии
Контрольная сумма cookbook\`ов стадии — это контрольная сумма всех [установленных файлов cookbook\`ов](#Установка-стадии-cookbookа) для данной стадии.

#### Дерево cookbooks
Дерево cookbooks — это результат выполнения berks vendor в chef приложении.

#### Shell проект
Shell проект — это проект, для сборки которого используются команды bash, заданные в [Dappfile](#dappfile).

#### Директория сборки проекта
Директория сборки проекта (build dir) — это директория для хранения кэша и производных данных сборки [проекта](#Проект).
* Используется для:
  * Хранения [дерева cookbooks](#Дерево-cookbooks).
    * Этот результат при возможности кэшируется между сборками.
  * Кэширования [контрольных сумм cookbook\`ов стадий](#Контрольная-сумма-cookbookов-стадии).
  * Хранения файлов блокировок проекта.
* Путь к build директории.
  * По умолчанию: \<[директория проекта](#Директория-проекта)\>/.dapps_build.
  * Переопределяется параметром --build-dir.
* Необходимо указывать одну и ту же build директорию для каждой из вызываемых команд dapp (или использовать путь к директории по умолчанию).
* При указании build директории необходимо иметь в виду, что для каждого [проекта](#Проект) нужна отдельная директория.
* Одна и та же build директория используется всеми [приложениями](#Приложение), описанными в данном [проекте](#Проект).

#### Временная директория приложения
Временная директория приложения (tmp dir) — это временная рабочая директория сборщика [приложения](#Приложение).
* Создается для каждого [приложения](#Приложение) во время сборки.
* Используется для хранения:
  * Склонированных remote git репозиториев.
  * [Стадий cookbook\`ов](#Стадия-cookbookа), [установленных](#Установка-стадии-cookbookа) для сборки стадии.
* Путь к директории: /tmp/dapp-\<date\>-\<random\>.
* Удаление директории происходит автоматически при ожидаемом успешном или не успешном завершении работы dapp.

#### Режим разработки
* Включается опцией командной строки --dev.
* Используется для:
  * включения создания файла Berksfile.lock в случае его отсутствия (см. [cookbook приложения](#cookbook-приложения)).

#### Блокировка ресурса
Блокировки обеспечивают корректную работу команд dapp при их параллельном запуске в рамках одного проекта на одном сервере.
* Блокировки реализованы с использованием механизма файловых блокировок ОС.
* Файлы блокировок хранятся в директории \<[build dir](#Директория-сборки)\>/locks.

#### Внешние директории сборки
Директории, которые монтируются в docker-контейнеры в процессе сборки [приложения](#Приложение).

* Используются для уменьшения размера [кэша приложения](#Кэш-приложения) и конечного образа приложения соответственно.
* Ожидаются директории с данными порождаемыми в процессе сборки, которые:
  * не требуются конечному пользователю;
  * занимают значительный размер.
* Для добавления используются директивы [tmp_dir](#tmp_dir) и [build_dir](#build_dir), которые различаются местом размещения внешних директорий. 

#### Стадии

| Имя                               | Краткое описание 					          | Зависимость от директив                            |
| --------------------------------- | ----------------------------------- | -------------------------------------------------- |
| from                              | Выбор окружения  					          | docker.from 			   						                   |
| before_install                    | Установка софта инфраструктуры      | shell.before_install / chef.module, chef.recipe    |
| before_install_artifact           | Наложение артефактов 				        | artifact (с before: :install) 			   		         |
| git_artifact_archive              | Наложение git-артефактов            | git_artifact.local`` и git_artifact.remote 		     |
| git_artifact_pre_install_patch    | Наложение патчей git-артефактов 	  | git_artifact.local и git_artifact.remote           |
| install                           | Установка софта приложения          | shell.install / chef.module, chef.recipe           |
| git_artifact_post_install_patch   | Наложение патчей git-артефактов     | git_artifact.local и git_artifact.remote           |
| after_install_artifact            | Наложение артефактов                | artifact (с after: :install)               		     |
| before_setup                      | Настройка софта инфраструктуры      | shell.before_setup / chef.module, chef.recipe      |
| before_setup_artifact             | Наложение артефактов                | artifact (с before: :setup)                		     |
| git_artifact_pre_setup_patch      | Наложение патчей git-артефактов     | git_artifact.local и git_artifact.remote           |
| setup                             | Развёртывание приложения            | shell.setup / chef.module, chef.recipe             |
| chef_cookbooks                    | Установка cookbook\`ов              | -             		       						               |
| git_artifact_post_setup_patch     | Наложение патчей git-артефактов     | git_artifact.local и git_artifact.remote           |
| after_setup_artifact              | Наложение артефактов                | artifact (с after: :setup)            	   		     |
| git_artifact_latest_patch         | Наложение патчей git-артефактов     | git_artifact.local и git_artifact.remote           |
| docker_instructions               | Применение докерфайловых инструкций | docker.cmd, docker.env, docker.entrypoint, docker.expose, docker.label, docker.onbuild, docker.user, docker.volume, docker.workdir |
| git_artifact_artifact_patch       | Наложение патчей git-артефактов     | git_artifact.local и git_artifact.remote           |
| build_artifact                    | Сборка артефакта                    | shell.build_artifact / chef.module, chef.recipe    |

##### Особенности
* Сигнатуры стадий git_artifact_pre_install_patch, git_artifact_post_install_patch, git_artifact_pre_setup_patch, помимо обратной зависимости от сигнатур стадий, имеют зависимость от install, before_setup и setup сигнатур соответственно.
  * К примеру: изменения зависимостей стадии install приведёт к пересборке стадии git_artifact_pre_install_patch.
* Сигнатура стадии git_artifact_post_setup_patch зависит от размера патчей git-артефактов и будет пересобрана, если их сумма превысит лимит (10 MB).

##### from
##### before install
##### before install artifact
##### git artifact archive
##### Группа install
###### git artifact pre install patch
###### install
###### git artifact post install patch
##### after install artifact
##### before setup
##### before setup artifact
##### Группа setup
###### git artifact pre setup patch
###### setup
###### chef cookbooks
###### git artifact post setup patch
##### after setup artifact
##### git artifact latest patch
##### docker instructions
##### git_artifact_artifact_patch
##### build_artifact

### Dappfile

#### Основное

##### name \<name\>
Базовое имя для собираемых docker image\`ей: \<базовое имя\>-dappstage:\<signature\>.

Опционально, по умолчанию определяется исходя из имени директории, в которой находится Dappfile.

##### install\_depends\_on \<glob\>[,\<glob\>, \<glob\>, ...]
Список файлов зависимостей для стадии install.

* При изменении содержимого указанных файлов, произойдет пересборка стадии install.
* Учитывается лишь содержимое файлов и порядок в котором они указаны (имена файлов не учитываются).
* Поддерживаются glob-паттерны.
* Директории игнорируются.

##### setup\_depends\_on \<glob\>[,\<glob\>, \<glob\>, ...]
Список файлов зависимостей для стадии setup.

* При изменении содержимого указанных файлов, произойдет пересборка стадии setup.
* Учитывается лишь содержимое файлов и порядок в котором они указаны (имена файлов не учитываются).
* Поддерживаются glob-паттерны.
* Директории игнорируются.

##### builder \<builder\>
Тип сборки: **:chef** или **:shell**.
* Опционально, по умолчанию будет выбран тот сборщик, который будет использован первым (см. [Chef](#chef), [Shell](#shell)).
* При определении типа сборки, сборщик другого типа сбрасывается.
* При смене сборщика, необходимо явно изменять тип сборки.
* Пример:
  * Собирать приложение X с **:chef** сборщиком, а Y c **:shell**:
  ```ruby
  app 'X' do
    chef.module 'a', 'b'
  end
  app 'Y' do
    shell.before_install 'apt-get install service'
  end
  ```
  * Собирать приложения X-Y и Z с **:chef** сборщиком, а X-V c **:shell**:
  ```ruby
  chef.module 'a', 'b'

  app 'X' do
    builder :shell

    app 'Y' do
      builder :chef
      chef.module 'c'
    end

    app 'V' do
      shell.install 'application install'
    end
  end
  app 'Z'
  ```

##### tmp_dir \<dir\>\[, dir\]
Определяет [внешние директории сборки](#Внешние-директории-сборки) для переданных путей **dir**, которые будут размещены во [временной директории приложения](#Временная-директория-приложения).

##### build_dir \<dir\>\[, dir\]
Определяет [внешние директории сборки](#Внешние-директории-сборки) для переданных путей **dir**, которые будут размещены в [директории сборки проекта](#Директория-сборки-проекта).

##### app \<app\>\[, &blk\]
Определяет приложение <app> для сборки.

* Опционально, по умолчанию будет использоваться приложение с базовым именем (см. name \<name\>).
* Можно определять несколько приложений в одном Dappfile.
* При использовании блока создается новый контекст.
  * Наследуются все настройки родительского контекста.
  * Можно дополнять или переопределять настройки родительского контекста.
  * Можно использовать директиву app внутри нового контекста.
* При использовании вложенных вызовов директивы, будут использоваться только приложения указанные последними в иерархии. Другими словами, в описанном дереве приложений будут использованы только листья.
  * Правила именования вложенных приложений: \<app\>[-\<subapp\>-\<subsubapp\>...]
* Примеры:
  * Собирать приложения X и Y:
  ```ruby
  app 'X'
  app 'Y'
  ```
  * Собирать приложения X, Y-Z и Y-K-M:
  ```ruby
  app 'X'
  app 'Y' do
    app 'Z'

    app 'K' do
      app 'M'
    end
  end
  ```

#### Артефакты

* Пути добавления не должны пересекаться между артефактами.
* Изменение любого параметра артефакта ведёт к смене сигнатур, пересборке связанных [стадий](#Стадии) приложения.
* Приложение может содержать любое количество артефактов.

##### Опциональные параметры артефакта

Все артефакты имеют следующие опциональные параметры: 

* cwd: поменять рабочую директорию.
* paths: добавить только указанные пути.
* exclude_paths: игнорировать указанные пути.  
* owner: определить пользователя.
* group: определить группу.

##### artifact \<where_to_add\>, before: \<before\>, after: \<after\> \[, [опциональные параметры артефакта](#Опциональные-параметры-артефакта)\]
Позволяет определять [артефакты](#Артефакт).

* Обязательный параметр **where_to_add** одновременно определяет директории артефакта, с необходимыми ресурсами, и связанных приложений, в которые будут копироваться эти ресурсы.
* Обязательно использование хотя бы одного из параметров **before** или **after**, где:
  * параметр определяет порядок применения артефакта (до или после);
  * значение определяет стадию (install или setup).

##### artifact\_depends\_on \<glob\>[,\<glob\>, \<glob\>, ...]
Список файлов зависимостей для стадии [build_artifact](#build_artifact) [артефакта](#Артефакт).

* При изменении содержимого указанных файлов, произойдет пересборка стадии build_artifact.
* Учитывается лишь содержимое файлов и порядок в котором они указаны (имена файлов не учитываются).
* Поддерживаются glob-паттерны.
* Директории игнорируются.

##### git_artifact.local \<where_to_add\>, \[, [опциональные параметры артефакта](#Опциональные-параметры-артефакта)\]
Позволяет определять локальный git-артефакт.

* Обязательный параметр **where_to_add** определяет директорию хранения артефакта.

##### git_artifact.remote \<url\>, \<where_to_add\>, \[, [опциональные параметры артефакта](#Опциональные-параметры-артефакта)\]
Позволяет определять удалённый git-артефакт.

* Обязательный параметр **where_to_add** определяет директорию хранения артефакта.
* Обязательный параметр **url** определяет адрес удалённого git-репозитория.

#### Docker

##### docker.from \<image\>[, cache_version: \<cache_version\>]
Определить окружение приложения **\<image\>** (см. [Стадия from](#from)).

* **\<image\>** имеет следующий формат 'REPOSITORY:TAG'.
* Опциональный параметр **\<cache_version\>** участвует в формировании сигнатуры стадии.

##### docker.cmd \<cmd\>[, \<cmd\> ...]
Применить dockerfile инструкцию CMD (см. [CMD](https://docs.docker.com/engine/reference/builder/#/cmd "Docker reference")).

##### docker.env \<env_name\>: \<env_value\>[, \<env_name\>: \<env_value\> ...]
Применить dockerfile инструкцию ENV (см. [ENV](https://docs.docker.com/engine/reference/builder/#/env "Docker reference")).

##### docker.entrypoint \<cmd\>[, \<arg\> ...]
Применить dockerfile инструкцию ENTRYPOINT (см. [ENTRYPOINT](https://docs.docker.com/engine/reference/builder/#/entrypoint "Docker reference")).

##### docker.expose \<expose\>[, \<expose\> ...]
Применить dockerfile инструкцию EXPOSE (см. [EXPOSE](https://docs.docker.com/engine/reference/builder/#/expose "Docker reference")).

##### docker.label \<label_key\>: \<label_value\>[, \<label_key\>: \<label_value\> ...]
Применить dockerfile инструкцию LABEL (см. [LABEL](https://docs.docker.com/engine/reference/builder/#/label "Docker reference")).

##### docker.onbuild \<cmd\>[, \<cmd\> ...]
Применить dockerfile инструкцию ONBUILD (см. [ONBUILD](https://docs.docker.com/engine/reference/builder/#/onbuild "Docker reference")).

##### docker.user \<user\>
Применить dockerfile инструкцию USER (см. [USER](https://docs.docker.com/engine/reference/builder/#/user "Docker reference")).

##### docker.volume \<volume\>[, \<volume\> ...]
Применить dockerfile инструкцию VOLUME (см. [VOLUME](https://docs.docker.com/engine/reference/builder/#/volume "Docker reference")).

##### docker.workdir \<path\>
Применить dockerfile инструкцию WORKDIR (см. [WORKDIR](https://docs.docker.com/engine/reference/builder/#/workdir "Docker reference")).

#### Shell

##### shell.before_install, shell.before_setup, shell.install, shell.setup [, cache_version: \<cache_version\>]

Директивы позволяют добавить bash-комманды для выполнения на соответствующих стадиях [shell приложения](#Shell-приложение).

* Опциональный параметр **\<cache_version\>** участвует в формировании сигнатуры стадии.

##### shell.build_artifact [, cache_version: \<cache_version\>]

Позволяет добавить bash-комманды для выполнения на соответствующей стадии [артефакта](#Артефакт).

* Опциональный параметр **\<cache_version\>** участвует в формировании сигнатуры стадии.

##### shell.reset_before_install, shell.reset_before_setup, shell.reset_setup, shell.reset_install, shell.reset_build_artifact

Позволяет сбросить объявленные ранее bash-комманды соответсвующей стадии.

##### shell.reset_all

Позволяет сбросить все объявленные ранее bash-комманды стадий.

#### Chef

##### chef.module \<mod\>[, \<mod\>, \<mod\> ...]
Включить переданные [модули](#mdapp-модуль) для chef builder в данном контексте.

* Для каждого переданного модуля может существовать по одному рецепту на каждую из стадий.
* При отсутствии файла рецепта в runlist для данной стадии используется пустой рецепт \<mod\>::void.

Подробнее см.: [mdapp модуль](#mdapp-модуль) и [установка стадии cookbook\`а](#Установка-стадии-cookbookа).

##### chef.skip_module \<mod\>[, \<mod\>, \<mod\> ...]
Выключить переданные модули для chef builder в данном контексте.

##### chef.reset_modules
Выключить все модули для chef builder в данном контексте.

##### chef.recipe \<recipe\>[, \<recipe\>, \<recipe\> ...]
Включить переданные рецепты из [приложения](#cookbook-приложения) для chef builder в данном контексте.

* Для каждого преданного рецепта может существовать файл рецепта в проекте на каждую из стадий.
* При отсутствии хотя бы одного файла рецепта из включенных, в runlist для данной стадии используется пустой рецепт \<projectname\>::void.
* Порядок вызова рецептов в runlist совпадает порядком их описания в конфиге.

Подробнее см.: [cookbook приложения](#cookbook-приложения) и [установка стадии cookbook\`а](#Установка-стадии-cookbookа).

##### chef.remove_recipe \<recipe\>[, \<recipe\>, \<recipe\> ...]
Выключить переданные рецепты из [приложения](#cookbook-приложения) для chef builder в данном контексте.

##### chef.reset_recipes
Выключить все рецепты из проекта для chef builder в данном контексте.

#### chef.attributes
Хэш атрибутов, доступных на всех стадиях сборки, для chef builder в данном контексте.

* Вложенные хэши создаются автоматически при первом обращении к методу доступа по ключу (см. пример).
* Пример использования:
```ruby
chef.attributes['mdapp-test']['nginx']['package_name'] = 'nginx-common'
chef.attributes['mdapp-test']['nginx']['package_version'] = '1.4.6-1ubuntu3.5'

app 'X' do
  chef.attributes['mdapp-test']['nginx']['package_version'] = '1.4.6-1ubuntu3'
end
```

См.: [установка стадии cookbook\`а](#Установка-стадии-cookbookа).

#### chef.\<стадия\>_attributes
Хэш атрибутов, доступных на стадии сборки, для chef builder в данном контексте.

См.: [установка стадии cookbook\`а](#Установка-стадии-cookbookа).

#### chef.reset_attributes
Выключить атрибуты, доступные на всех стадиях сборки, для chef builder в данном контексте.

#### chef.reset_\<стадия\>_attributes
Выключить атрибуты, доступные на стадии сборки, для chef builder в данном контексте.

#### chef.reset_all_attributes
Выключить все атрибуты для chef builder в данном контексте.

#### chef.reset_all
Выключить все рецепты из [приложения](#cookbook-приложения), все модули для chef builder и все атрибуты в данном контексте.

### Команды

#### dapp build
Собрать приложения, удовлетворяющие хотя бы одному из **APPS PATTERN**-ов (по умолчанию *).

```
dapp build [options] [APPS PATTERN ...]
```

##### Опции среды сборки

###### --dir PATH
Определяет директорию, которая используется при поиске одного или нескольких **Dappfile**.

По умолчанию поиск ведётся в текущей папке пользователя.

###### --build-dir PATH
Переопределяет директорию хранения кэша, который может использоваться между сборками.

###### --tmp-dir-prefix PREFIX
Переопределяет префикс временной директории, файлы которой используются только во время сборки.

##### Опции логирования

###### --dry-run
Позволяет запустить сборщик вхолостую и посмотреть процесс сборки.

###### --verbose
Подробный вывод.

###### --color MODE
Отвечает за регулирование цвета при выводе в терминал.

Существует несколько режимов (**MODE**): **on**, **of**, **auto**.

По умолчанию используется **auto**, который окрашивает вывод, если вывод производится непосредственно в терминал.

###### --time
Добавляет время каждому событию лога.

##### Опции интроспекции
Позволяют поработать с образом на определённом этапе сборки.

###### --introspect-stage STAGE
После успешного прохождения стадии **STAGE**.

###### --introspect-before-error
Перед выполением команд несобравшейся стадии.

###### --introspect-error
После завершения команд стадии с ошибкой.

##### Примеры
* Сборка в текущей директории:
```bash
$ dapp build
```
* Сборка приложений из соседней директории:
```bash
$ dapp build --dir ../project
```
* Запуск вхолостую с подробным выводом процесса сборки:
```bash
$ dapp build --dry-run --verbose
```
* Выполнить сборку, а в случае ошибки, предоставить образ для тестирования:
```bash
$ dapp build --introspect-error
```

#### dapp push
Выкатить собранное приложение в репозиторий, в следующем формате **REPO**:**ИМЯ ПРИЛОЖЕНИЯ**-**TAG**.

```
dapp push [options] [APPS PATTERN ...] REPO
```

##### --with-stages
Также выкатить кэш стадий.

##### --force
Позволяет перезаписывать существующие образы.

##### Опции тегирования
Отвечают за тег(и), с которыми выкатывается приложение.

Могут быть использованы совместно и по несколько раз.

В случае отсутствия, используется тег **latest**.

###### --tag TAG
Добавляет произвольный тег **TAG**.

###### --tag-branch
Добавляет тег с именем ветки сборки.

###### --tag-commit
Добавляет тег с коммитом сборки.

###### --tag-build-id
Добавляет тег с идентификатором сборки (CI).

###### --tag-ci
Добавляет теги, взятые из переменных окружения CI систем.

##### Примеры
* Выкатить все приложения в репозиторий localhost:5000/test и тегом latest:
```bash
$ dapp push localhost:5000/test
```
* Посмотреть, какие образы могут быть добавлены в репозиторий:
```bash
$ dapp push localhost:5000/test --tag yellow --tag-branch --dry-run
backend
  localhost:5000/test:backend-yellow
  localhost:5000/test:backend-master
frontend
  localhost:5000/test:frontend-yellow
  localhost:5000/test:frontend-0.2
```

#### dapp spush
Выкатить собранное приложение в репозиторий, в следующем формате **REPO**:**TAG**.

```
dapp spush [options] [APP PATTERN] REPO
```

Опции такие же как у **dapp push**.

##### Примеры
* Выкатить приложение **app** в репозиторий test, именем myapp и тегом latest:
```bash
$ dapp spush app localhost:5000/test
```
* Выкатить приложение с произвольными тегами:
```bash
$ dapp spush app localhost:5000/test --tag 1 --tag test
```
* Посмотреть, какие образы могут быть добавлены в репозиторий:
```bash
$ dapp spush app localhost:5000/test --tag-commit --tag-branch --dry-run
localhost:5000/test:2c622c16c39d4938dcdf7f5c08f7ed4efa8384c4
localhost:5000/test:master
```

#### dapp stages push
Выкатить кэш собранных приложений в репозиторий.

```
dapp stages push [options] [APP PATTERN] REPO
```

##### Примеры
* Выкатить кэш приложений проекта в репозиторий localhost:5000/test:
```bash
$ dapp stages push localhost:5000/test
```
* Посмотреть, какие образы могут быть добавлены в репозиторий:
```bash
$ dapp stages push localhost:5000/test --dry-run
backend
  localhost:5000/test:dappstage-be032ed31bd96506d0ed550fa914017452b553c7f1ecbb136216b2dd2d3d1623
  localhost:5000/test:dappstage-2183f7db73687e727d9841594e30d8cb200312290a0a967ef214fe3771224ee2
  localhost:5000/test:dappstage-f7d4c5c420f29b7b419f070ca45f899d2c65227bde1944f7d340d6e37087c68d
  localhost:5000/test:dappstage-256f03ccf980b805d0e5944ab83a28c2374fbb69ef62b8d2a52f32adef91692f
  localhost:5000/test:dappstage-31ed444b92690451e7fa889a137ffc4c3d4a128cb5a7e9578414abf107af13ee
  localhost:5000/test:dappstage-02b636d9316012880e40da44ee5da3f1067cedd66caa3bf89572716cd1f894da
frontend
  localhost:5000/test:dappstage-192c0e9d588a51747ce757e61be13acb4802dc2cdefbeec53ca254d014560d46
  localhost:5000/test:dappstage-427b999000024f9268a46b889d66dae999efbfe04037fb6fc0b1cd7ebb4600b0
  localhost:5000/test:dappstage-07fe13aec1e9ce0fe2d2890af4e4f81aaa984c89a2b91fbd0e164468a1394d46
  localhost:5000/test:dappstage-ba95ec8a00638ddac413a13e303715dd2c93b80295c832af440c04a46f3e8555
  localhost:5000/test:dappstage-02b636d9316012880e40da44ee5da3f1067cedd66caa3bf89572716cd1f894da
```

#### dapp stages pull
Импортировать необходимый кэш приложений проекта, если он присутствует в репозитории **REPO**.

Если не указана опция **--all**, импорт будет выполнен до первого найденного кэша для каждого приложения.

```
dapp stages pull [options] [APP PATTERN] REPO
```

##### --all
Попробовать импортировать весь необходимый кэш.

##### Примеры
* Импортировать кэш приложений проекта из репозитория localhost:5000/test:
```bash
$ dapp stages pull localhost:5000/test
```
* Посмотреть, поиск каких образов в репозитории localhost:5000/test может быть выполен:
```bash
$ dapp stages pull localhost:5000/test --all --dry-run
backend
  localhost:5000/test:dappstage-be032ed31bd96506d0ed550fa914017452b553c7f1ecbb136216b2dd2d3d1623
  localhost:5000/test:dappstage-2183f7db73687e727d9841594e30d8cb200312290a0a967ef214fe3771224ee2
  localhost:5000/test:dappstage-f7d4c5c420f29b7b419f070ca45f899d2c65227bde1944f7d340d6e37087c68d
  localhost:5000/test:dappstage-256f03ccf980b805d0e5944ab83a28c2374fbb69ef62b8d2a52f32adef91692f
  localhost:5000/test:dappstage-31ed444b92690451e7fa889a137ffc4c3d4a128cb5a7e9578414abf107af13ee
  localhost:5000/test:dappstage-02b636d9316012880e40da44ee5da3f1067cedd66caa3bf89572716cd1f894da
frontend
  localhost:5000/test:dappstage-192c0e9d588a51747ce757e61be13acb4802dc2cdefbeec53ca254d014560d46
  localhost:5000/test:dappstage-427b999000024f9268a46b889d66dae999efbfe04037fb6fc0b1cd7ebb4600b0
  localhost:5000/test:dappstage-07fe13aec1e9ce0fe2d2890af4e4f81aaa984c89a2b91fbd0e164468a1394d46
  localhost:5000/test:dappstage-ba95ec8a00638ddac413a13e303715dd2c93b80295c832af440c04a46f3e8555
  localhost:5000/test:dappstage-02b636d9316012880e40da44ee5da3f1067cedd66caa3bf89572716cd1f894da
```

#### dapp list
Вывести список приложений.

```
dapp list [options] [APPS PATTERN ...]
```

#### dapp run
Запустить собранное приложение с докерными аргументами **DOCKER ARGS**.

```
dapp run [options] [APPS PATTERN...] [DOCKER ARGS]
```

##### [DOCKER ARGS]
Может содержать докерные опции и/или команду.

Перед командой необходимо использовать группу символов ' -- '.

##### Примеры
* Запустить приложение с опциями:
```bash
$ dapp run -ti --rm
```
* Запустить с опциями и командами:
```bash
$ dapp run -ti --rm -- bash -ec true
```
* Запустить, передав только команды:
```bash
$ dapp run -- bash -ec true
```
* Посмотреть, что может быть запущено:
```bash
$ dapp run app -ti --rm -- bash -ec true
docker run -ti --rm app-dappstage:ea5ec7543c809ec7e9fe28181edfcb2ee6f48efaa680f67bf23a0fc0057ea54c bash -ec true
```

#### dapp stages cleanup local
Удалить неактуальный локальный кэш приложений проекта, опираясь на приложения в репозитории **REPO**.

```
dapp stages cleanup local [options] [APPS PATTERN ...] REPO
```

##### --improper-cache-version
Удалить устаревший кэш приложений проекта.

##### Примеры
* Удалить неактуальный кэш приложений:
```bash
$ dapp stages cleanup local localhost:5000/test --improper-cache-version
```

#### dapp stages cleanup repo
Удалить неиспользуемый кэш приложений в репозитории **REPO**.

```
dapp stages cleanup repo [options] [APPS PATTERN ...] REPO
```

##### --improper-cache-version
Удалить устаревший кэш приложений проекта.

##### Примеры
* Удалить неактуальный кэш приложений в репозитории localhost:5000/test:
```bash
$ dapp stages cleanup repo localhost:5000/test
```

#### dapp stages flush local
Удалить кэш приложений проекта.

```
dapp stages flush local [options] [APPS PATTERN ...]
```

##### Примеры
* Удалить кэш приложений:
```bash
$ dapp stages flush local
```

#### dapp stages flush repo
Удалить приложения и кэш приложений проекта в репозитории **REPO**.

```
dapp stages flush repo [options] [APPS PATTERN ...] REPO
```

##### Примеры
* Удалить весь кэш приложений в репозитории localhost:5000/test:
```bash
$ dapp stages flush repo localhost:5000/test
```

#### dapp cleanup
Убраться в системе после некорректного завершения работы dapp, удалить нетеггированные docker-образы и docker-контейнеры проекта.

```
dapp cleanup [options] [APPS PATTERN ...]
```

##### Примеры
* Запустить:
```bash
$ dapp cleanup
```
* Посмотреть, какие команды могут быть выполнены:
```bash
$ dapp cleanup --dry-run
backend
  docker rm -f dd4ec7v33
  docker rmi ea5ec7543 c809ec7e9f ee6f48efa6
```

#### dapp mrproper
Очистить систему.

```
dapp mrproper [options]
```

##### --all
Удалить docker-образы и docker-контейнеры связанные с dapp.

##### --improper-cache-version-stages
Удалить устаревший кэш приложений.

##### Примеры
* Запустить очистку:
```bash
$ dapp mrproper --all
```
* Посмотреть, версия кэша каких образов устарела, какие команды могут быть выполнены:
```bash
$ dapp mrproper --improper-cache-version-stages --dry-run
mrproper
  proper cache
    docker rmi dappstage-dapp-test-project-services-stats:ba95ec8a00638ddac413a13e303715dd2c93b80295c832af440c04a46f3e8555 dappstage-dapp-test-project-services-stats:f53af70566ec23fb634800d159425da6e7e61937afa95e4ed8bf531f3503daa6
```
