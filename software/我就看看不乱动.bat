@echo off
start cmd /k "cd/d kafka_2.11-2.4.1 &&bin\windows\kafka-console-consumer.bat --bootstrap-server=127.0.0.1:9092 --topic=web_log --from-beginning"