# WB go advert bid checker

Проверяет цену ставки на рекламу в Wildberries по списку ключевых слов.

Пример запроса:
```
curl http://127.0.0.1/advert-bet-info -d '{
                                 "subject_id": 226, "target_place": 5, "validate_subject_id_priority": true, "advert_id": 3743672,
                                 "keywords": [
                                     "нитриловые перчатки s",
                                     "нитриловые перчатки м",
                                     "перчатки латексные",
                                     "перчатки резиновые",
                                     "нитриловые перчатки",
                                     "перчатки",
                                     "одноразовые перчатки",
                                     "резиновые перчатки"
                                 ]
                             }'
```
* subject_id - категория
* target_place - порядковый номер рекламного места
* validate_subject_id_priority - проверять ли что по ключевому слову запрашиваемая категория самая приоритетная. Все ключевые слова по которым есть боллее приоритетрые категории будут собраны в ответе в атрибуте warnings
* advert_id - id рекламной кампании, по нуму будет найдено текущее занимаемое место
* keywords - список ключевых слов


Ответ: 
```
{
    "bets": {
        "target_place": {
            "keyword": "перчатки",
            "subject_id": 226,
            "advert_id": 7985287,
            "place": 5,
            "bet": 300
        },
        "next_place": {
            "keyword": "перчатки",
            "subject_id": 226,
            "advert_id": 7986352,
            "place": 6,
            "bet": 300
        },
        "my_place": 8
    },
    "warnings": [
        {
            "keyword": "одноразовые перчатки",
            "priority_subject_id": 2652
        },
        {
            "keyword": "нитриловые перчатки",
            "priority_subject_id": 2652
        }
    ]
}
```