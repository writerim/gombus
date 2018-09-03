go get http://code.nag.how/writer/gombus.git


#### API

Создание нового объекта  
**New( ) ( *gombus.Mbus )**


Присвоить серийный номер   
**(*gombus.Mbus) SetSerial( ) ( error )**

Получить запрос для открытия соединения  
**(*gombus.Mbus) OpenCmd( ) ( []byte , error )**

Получить запрос для считывания данных (длинный запрос)  
**(*gombus.Mbus) ReadLongCmd( ) ( []byte , error )**
 
Парсер длянного ответа  
**(*gombus.Mbus) ParseLongCmd( ) ( map[string]int , error )**