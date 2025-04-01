# news-feed-bot

### git-readme.

## База (Common commands)

#### 1. `git status` (тыкать в любом непонятном случае)

#### 2. `git add .` (добавлять файлы в `stage`)

#### 3. `git commit -m ""` (`commit` - запись)

#### 4. `git log` (посмотреть кто делал комит и его хеш)

#### 5. `git push origin main` (отправить на удаленный репозиторий. `origin` - вместо ссылки на реп, `main` - название ветки)

#### 6. `git pull` (подтягивает изменения с удаленного репозитория)

## Изменение файлов (File changing)

#### 1. `git reset [file_name]` (убрать некоторые файлы из промежуточной области для одного файла)

#### 2. `git diff` (показать изменения в файлах)

#### 3. `git reset --hard` (убрать изменения из ВСЕХ файлов)

## Создание веток (Branch creation)

#### 1. `git branch` (просмотр всех веток и на какой ты)

#### 2. `git branch [branch_name]` (создание ветки)

#### 3. `git checkout [branch_name]` (переключение между ветками)

#### 4. `git checkout -b [branch_name]` (создание ветки и переключение на нее)

## Удаление веток (Branch deletion)

#### 1. `git branch -d [branch_name]` (флаг `-d` это удаление ветки)

## Слияние веток (Branch merging)

#### 1. `git merge [branch_name]` (слияние веток. `[branch_name]` - ветка, из которой берем изменения (находимся в той ветке, в которую закидываем изменения))

#### 2. GitHub -> project -> Pull Request -> New pull request -> Create pull request -> Add title + comment -> Create pull request

#### 3. add comment -> start review -> finish review

#### 4. resolve conversation -> merge pull request

#### 5. goood job

## Разделы в pull request

#### 1. conversation

#### 2. files changed

## Решение конфликтов (Conflict resolving)

#### 1. `git merge` -> Conflict occured ->

<span style="color:white; font-weight:500; font-size:16px"><span style="color:green; font-weight:700; ">Зеленая полоска</span> - ( <span style="color:green; font-weight:700; "><<< CURRENT CHANGE</span> ) это где находишься,</span>

<span style="color:white; font-weight:500; font-size:16px "><span style="color:#ADD8E6; font-weight:700; ">Cиняя полоска</span> - ( <span style="color:#ADD8E6; font-weight:700; ">>>> INCOMING CHANGE</span> ) другая ветка, в которой тоже меняли этот файл.</span>

#### -> Выбрать надпись сверху:

- `Accept Current Change` - Принять изменения из CURRENT CHANGE
- `Accept Incoming Change` - Принять изменения из INCOMING CHANGE
- `Accept Both Change` - Принять изменения из Обоих вариантов

#### -> `git add .` -> `git commit -m [message]` -> `git log`

## Работа с GitFlow

#### 1. Создать репозиторий на GitHub и клонировать на компьютер

#### 2. Создать ветку разработки `development` от главной ветки

#### 3. Создать от ветки `development` **feature-ветки** и мержить **feature-ветки** в `development`, когда фичи готовы.

#### 4. Создание ветки `release/0.1.0` от `development`

#### 5. Когда ветка `release/0.1.0` закончена, то она мержится в `development` и `main` и затем **удаляется**

#### 6. Если в ветке `main` обнаруживается ошибка, то создается `hotfix-ветка`

#### 7. Когда работа над `hotfix-веткой` завершается, ее нужно мержить в `development` и `main`, а затем **удалить**

## Доступ к репозиторию по SSH

#### 1. Создать приватный репозиторий на github

#### 2. `ssh-keygen -o` -> enter passphrase -> cat + [public key] (скопировать путь к файлу без расширения)

#### 3. Скопировать ssh-ключ -> зайти на github -> settings -> ssh and gpg keys -> new ssh key -> вставить ключ и ОСОЗНАНОЕ имя ключа -> add ssh key

#### 4. В созданном репозитории скопировать ssh url

#### 5. В терминале `ssh-add + путь к файлу` `(пример /User/sofiya/.ssh/test-ssh)` -> ввести пароль

#### 6. `git clone [ssh-url]`

#### 7. Открыть созданную папку