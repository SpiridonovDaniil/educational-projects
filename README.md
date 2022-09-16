# educational-projects
I am learning the golang programming language

Задача 1

Цель практической работы

Научиться: 

    писать микросервис и proxy,
    тестировать написанное приложение. 


Что нужно сделать

В прошлом домашнем задании вы писали приложение, которое принимает HTTP-запросы, создаёт пользователей, добавляет друзей и так далее. 
Давайте теперь приблизим наше приложение к реальному продукту. 

    Отрефакторьте приложение так, чтобы вы могли поднять две реплики данного приложения. 
    Используйте любую базу данных, чтобы сохранять информацию о пользователях, или можете сохранять информацию в файл, предварительно сереализуя в JSON. 
    Напишите proxy или используйте, например, nginx. 
    Протестируйте приложение.

В данном проекте реализован  HTTP-сервис, который принимает входящие соединения с JSON-данными и обрабатывает их следующим образом:
- создание пользователя;
- добавление друзей;
- удаление пользователя;
- возврат друзей пользователя;
- обновление возраста пользователя.

Также был реализован небольшой микросервис и proxy. Сервис был протестирован с использованием пакета: github.com/golang/mock/gomock.
Для хранения информации использовалась база данных mongodb. 
Для работы с которой использовался пакет: gopkg.in/mgo.v2. 
В проекте была реализована контейнеризация приложения с помощью docker и docker-compose.




Задание 2. Конвейер
Цели задания

    Научиться работать с каналами и горутинами.
    Понять, как должно происходить общение между потоками.

Что нужно сделать

Реализуйте паттерн-конвейер: 

    Программа принимает числа из стандартного ввода в бесконечном цикле и передаёт число в горутину.
    Квадрат: горутина высчитывает квадрат этого числа и передаёт в следующую горутину.
    Произведение: следующая горутина умножает квадрат числа на 2.
    При вводе «стоп» выполнение программы останавливается. 

Советы и рекомендации

Воспользуйтесь небуферизированными каналами и waitgroup.
Что оценивается

Ввод : 3

Квадрат : 9

Произведение : 18
Как отправить задание на проверку

Выполните задание в файле вашей среды разработки и пришлите ссылку на архив с вашим проектом через форму ниже.






Задание 3. Graceful shutdown
Цель задания

Научиться правильно останавливать приложения.
Что нужно сделать

В работе часто возникает потребность правильно останавливать приложения. Например, когда наш сервер обслуживает соединения, а нам хочется, чтобы все текущие соединения были обработаны и лишь потом произошло выключение сервиса. Для этого существует паттерн graceful shutdown. 

Напишите приложение, которое выводит квадраты натуральных чисел на экран, а после получения сигнала ^С обрабатывает этот сигнал, пишет «выхожу из программы» и выходит.
Советы и рекомендации

Для реализации данного паттерна воспользуйтесь каналами и оператором select с default-кейсом.
Что оценивается

Код выводит квадраты натуральных чисел на экран, после получения ^С происходит обработка сигнала и выход из программы.
Как отправить задание на проверку

Выполните задание в файле вашей среды разработки и пришлите ссылку на архив с вашим проектом через форму ниже.
