# atlant.io

- Fetch(URL) - запросить внешний CSV-файл со списком продуктов по внешнему адресу.CSV-файл имеет вид PRODUCT NAME;PRICE. Последняя цена каждого продукта должна быть сохранена в базе с датой запроса. Также нужно сохранять количество изменений цены продукта.

- List(<paging params>, <sorting params>) - получить постраничный список продуктов с их ценами, количеством изменений цены и датами их последнего обновления.Предусмотреть все варианты сортировки для реализации интерфейса в виде бесконечного скролла.

* Сервер должен быть запущен в 2+ экземплярах (каждый в своем Docker-контейнере) изакрыт балансировщиком, соответствующие конфигурации также должны бытьпредоставлены для тестовой среды.

# create upstream
curl -X POST http://localhost:8001/upstreams --data "name=price.v1.service"

# add two targets to the upstream
$ curl -X POST http://localhost:8001/upstreams/price.v1.service/targets \
    --data "target=api-1:50051"

$ curl -X POST http://localhost:8001/upstreams/price.v1.service/targets \
    --data "target=api-2:50051"

curl -XPOST http://localhost:8001/services/ \
  --data name=grpc \
  --data protocol=grpc \
  --data "host=price.v1.service" \
  --data port=15002

curl -XPOST http://localhost:8001/services/grpc/routes \
  --data protocols=grpc \
  --data name=catch-all \
  --data paths=/