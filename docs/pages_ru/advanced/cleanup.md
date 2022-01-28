---
title: Очистка
permalink: advanced/cleanup.html
---

При работе werf сохраняет данные как в container registry, так и на хосте, на котором запускается. 

В случае с container registry очистка требуется для удаления неактуальных образов с учётом установленных политик очистки, а также используемых образов в Kubernetes. 
На хосте все данные werf можно разделить на кеш, временные данные, которые остаются после запуска werf и более не требуются, а также локальные стадии в Docker, создаваемые при использовании werf без container registry.

Очистка container registry и хоста выполняется раздельно независимыми командами. 
При очистке container registry поддерживается работа только с данными конкретного проекта, а очистка хоста по умолчанию выполняется сразу для всех проектов. 

|                                                       | **Проект**                                       | **Все проекты**       |
|-------------------------------------------------------| :----------------------------------------------: | :-------------------: |
| Очистка неактуальных данных в **container registry**  | `werf cleanup --repo REPO`                       | -                     |
| Полная очистка **container registry**                 | `werf purge --repo REPO`                         | -                     |
| Очистка неактуальных данных **хоста**                 | `werf host cleanup --project-name PROJECT`__*__  | `werf host cleanup`   |
| Полная очистка **хоста**                              | `werf host purge --project-name PROJECT`__*__    | `werf host purge`     |

_, где __*__ указывает на неполную поддержку._ 

Следует отметить, что:
- **Очистку неактуальных данных** безопасно вызывать **в любой момент времени**, вручную или **автоматически**, с малыми рисками потерять критически важные данные, которые используются в production.
  - Более того werf без дополнительных настроек может производить очистку неактуальных данных хоста во время работы произвольных команд автоматически.
- **Полную очистку** данных не предполагается вызывать автоматически, только **вручную** и со знанием дела, потому что могут быть потеряны критически важные данные, которые используются в production.

## Очистка container registry

### Очистка неактуальных данных

Команда [**werf cleanup**]({{ "reference/cli/werf_cleanup.html" | true_relative_url }}) рассчитана на периодический запуск по расписанию. Удаление производится в соответствии с принятыми политиками очистки и является безопасной процедурой.

Алгоритм автоматически за пользователя выбирает образы, которые следует удалить. Алгоритм очистки может быть представлен следующими шагами:

- Получение необходимых данных из container registry.
- Подготовка списка образов, которые не должны быть удалены:
  - [Образы, которые используются в Kubernetes](#образы-в-kubernetes).
  - Образы, которые попадают под [пользовательские политики](#пользовательские-политики) при [сканировании истории git](#сканирование-истории-git).
  - Свежие образы, которые собирались за определённый период времени (регулируется опцией `--keep-stages-built-within-last-n-hours`, по умолчанию за последние два часа).
  - Связанные образы для получившегося на предыдущих шагах списка.   
- Удаление оставшихся образов.

#### Образы в Kubernetes

werf подключается **ко всем кластерам** Kubernetes, описанным **во всех контекстах** конфигурации kubectl, и собирает имена образов для следующих типов объектов: `pod`, `deployment`, `replicaset`, `statefulset`, `daemonset`, `job`, `cronjob`, `replicationcontroller`.

Пользователь может регулировать поведение следующими параметрами (и связанными переменными окружения):
- `--kube-config`, `--kube-config-base64` для определения конфигурации kubectl (по умолчанию используется пользовательская конфигурация `~/.kube/config`).
- `--kube-context` для выполнения сканирования только в определённом контексте.
- `--scan-context-namespace-only` для сканирования только связанного с контекстом namespace (по умолчанию все).
- `--without-kube` для отключения сканирования Kubernetes.

Пока в кластере Kubernetes существует объект использующий образ, он никогда не удалится из container registry. Другими словами, если что-то было запущено в вашем кластере Kubernetes, то используемые образы ни при каких условиях не будут удалены при очистке.

#### Сканирование истории git

В основу алгоритма очистки ложится тот факт, что в container registry сохраняется информация о коммитах, на которых выполняется сборка (добавился, изменился или нет образ в container registry — не имеет значения). При каждой сборке сохраняется связка коммит, [дайджест стадии]({{ "internals/stages_and_storage.html#дайджест-стадии" | true_relative_url }}) и имя образа — для каждого `image` из `werf.yaml`.

При сборке очередного коммита _конечный образ_ может не измениться, тем не менее в container registry добавится запись с информацией о том, что сборка [дайджеста стадии]({{ "internals/stages_and_storage.html#дайджест-стадии" | true_relative_url }}) соответствует определённому коммиту.

Таким образом, обеспечивается связь [дайджеста стадии]({{ "internals/stages_and_storage.html#дайджест-стадии" | true_relative_url }}) с историей git и появляется возможность организации эффективной очистки неактуальных образов на основе состояния git и [выбранных политик](#пользовательские-политики). Алгоритм сканирует историю git и отбирает значимые образы.

Информация о коммитах является единственным источником правды при работе алгоритма, поэтому образы без подобной информации будут удалены.

#### Пользовательские политики

Используя [политики очистки]({{ "advanced/cleanup.html" | true_relative_url }}), `keepPolicies`, пользователь определяет образы, которые не должны удаляться при очистке. При отсутствии конфигурации в `werf.yaml` будет использован [набор политик по умолчанию]({{ "reference/werf_yaml.html#политики-по-умолчанию" | true_relative_url }}).

Стоит отметить, что алгоритм сканирует локальное состояние git репозитория и актуальность git-веток и git-тегов крайне важна. По умолчанию werf выполняет синхронизацию автоматически (поведение регулируется в `werf.yaml` директивой [gitWorktree.allowFetchOriginBranchesAndTags]({{ "reference/werf_yaml.html#git-worktree" | true_relative_url }})).

#### Особенность очистки собираемых образов

При очистке пользовательские политики применяются к набору образов для каждого `image` из `werf.yaml`. Очистка должна учитывать все используемые `image` и для этого не всегда достаточно набора из основной ветки git-репозитория (к примеру, при разработке в feature-ветке могут добавляться/удаляться `image`).

Для того, чтобы избежать удаления рабочих образов и стадий, при сборке в container registry добавляется имя собираемого образа. Такой подход позволяет отвязаться от набора имён в `werf.yaml` при выполнении периодической очистки и учитывать все когда-либо используемые образы.  

Набор команд `werf managed-images ls|add|rm` позволяет пользователю редактировать так называемый набор _managed images_ и явно удалять образы, которые более не должны участвовать в очистке и могут быть полностью удалены.

### Полная очистка

Команда [**werf purge**]({{ "reference/cli/werf_purge.html" | true_relative_url }}) используется для полного удаления образов из container registry. Команда не учитывает, используются образы в кластере Kubernetes или нет.

Команда работает только в рамках проекта и требует доступа к git-репозиторию проекта, содержащему werf.yaml.

## Очистка хоста

### Очистка неактуальных данных

Команда [**werf host cleanup**]({{ "reference/cli/werf_host_cleanup.html" | true_relative_url }}) очищает старые, неиспользуемые и неактуальные данные, а так же сокращает размер кеша для всех проектов на хосте, учитывая занятое место и настройки пользователя.

Алгоритм автоматически за пользователя решает какие данные будут удалены. Алгоритм очистки может быть представлен следующими шагами:

 - Оценивается используемое место на том томе, где располагаются данные локального docker server.
 - Если используемое место превышает выставленный порог (по умолчанию это 70% занятого места на томе — задаётся параметром), то вычисляется количество места, которое требуется освободить, чтобы снизить используемое место на томе до выставленного уровня (по умолчанию это порог в 70% минус 5% запаса — задаётся параметром).
 - Далее алгоритм начинает удалять наиболее давно используемые данные (по принципу LRU — least recently used) до тех пор пока не будет удалено достаточно данных, чтобы уровень занятого места на диске снизился до требуемого порога (по умолчанию это порог в 70% минус 5% запаса — задаётся параметром).

Какие данные могут быть удалены:
 - Git-архивы из локального кеша werf: `~/.werf/local_cache/git_archives`.
 - Git-патчи из локального кеша werf: `~/.werf/local_cache/git_patches`.
 - Git-репозитории из локального кеша werf: `~/.werf/local_cache/git_repos`.
 - Git-worktree из локального кеша werf: `~/.werf/local_cache/git_worktrees`.
 - Все docker образы, собранные версией v1.2, оставшиеся в локальном docker server.
 - Те docker образы, собранные версией v1.1, которые хранятся в `--stages-storage=REPO`.
     - Образы от версии v1.1, которые хранятся в `--stages-storage=:local` не могут быть очищены данным алгоритмом, т.к. это первичное хранилище стадий, которые могут использоваться в production и других окружениях.

Следует отметить, что данный алгоритм в команде `werf host cleanup` применяется отдельно для тома где хранится локальный кеш werf `~/.werf/local_cache` и для тома, где хранятся данные локального docker server (обычно это `/var/lib/docker`). В том случае, если werf не может самостоятельно определить том, где хранятся реальные данные локального docker server, имеется возможность явно указать директорию данных локального docker server через параметр `--docker-server-storage-path=/var/lib/docker` (либо через переменную окружения `WERF_DOCKER_SERVER_STORAGE_PATH`).

По умолчанию при использовании werf очистка неактуальных данных хоста может выполняться автоматически в любой команде werf и нет никакой необходимости в дополнительных вызовах команды `werf host cleanup` вручную или в cron. Однако пользователь может выключить автоочистку неактуальных данных хоста с помощью параметра `--disable-auto-host-cleanup` (или переменной окружения `WERF_DISABLE_AUTO_HOST_CLEANUP`). В этом случае рекомендуется добавить команду `werf host cleanup` в cron, например следующим образом:

```shell
# /etc/cron.d/werf-host-cleanup
SHELL=/bin/bash
*/30 * * * * gitlab-runner source ~/.profile ; source $(trdl use werf 1.2 stable) ; werf host cleanup
```

По умолчанию без дополнительных параметров `werf host cleanup` будет чистить данные всех проектов на хосте. С параметром `--project-name PROJECT` команда может удалять только образы из локального docker-сервера. В данном режиме команда поддерживается частично.

### Полная очистка

Команда [**werf host purge**]({{ "reference/cli/werf_host_purge.html" | true_relative_url }}) имеет 2 режима работы: очистка данных одного проекта или очистка данных всех проектов.

**ВАЖНО.** По умолчанию без дополнительных параметров `werf host purge` удалит все следы werf на хосте: образы, стадии, кеш и другие данные (служебные папки, временные файлы) всех проектов. Команда обеспечивает максимальную степень очистки.

С параметром `--project-name PROJECT` команда удалит образы из локального docker server, связанные с данным проектом. В данном режиме команда поддерживается частично, не будут удалены образы в локальном docker server, связанные с удалённым хранилищем образов в container registry (например, локальные образы, оставшиеся от `werf converge --repo REPO`). Для удаления таких образов можно использовать команду очистки неактуальных данных хоста в режиме очистки данных всех проектов: `werf host cleanup`.