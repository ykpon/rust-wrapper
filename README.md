# README.md
- en [Русский](README.md)
- ru [English](README.en.md)

# Инструмент управления сервером Rust

Удобный инструмент для управления сервером игры Rust. Цель - предоставить возможность управления сервером из консоли Linux.

## Основные функции

1. **Управление сервером:** позволяет пользователям отправлять команды напрямую в консоль сервера Rust через протокол RCON. Так как сам по себе сервер Rust не предоставляет такой возможности из-за особенностей реализации игрового движка Unity, данная программа реализует этот функционал через RCON интерфейс.
2. **Логирование:** берет на себя функции вывода потока вывода консоли сервера, выводя все обычные сообщения и ответы RCON в том числе.
3. **Совместимость с Docker:** разработан в том числе для использования сервера Rust в Docker.

## Переменные окружения

Использует следующие переменные окружения для конфигурации:

- `RCON_IP`: IP-адрес сервера Rust. Если не указано, по умолчанию используется `127.0.0.1`.
- `RCON_PORT`: RCON порт сервера Rust. Если не указано, по умолчанию используется `28018`.
- `RCON_PASS`: RCON пароль сервера Rust. Значение по умолчанию отсутствует.

## Использование

После установки соответствующих переменных окружения вы можете запустить программу с исполняемым файлом вашего сервера Rust и любыми параметрами сервера.

## Примечание

Инструмент находится в бета-состоянии, и как таковой, он может быть не полностью функционален. Поэтому его не рекомендуется использовать в производственной среде без тщательного тестирования. Пожалуйста, не стесняйтесь вносить свой вклад, создавая проблемы или пул-реквесты.