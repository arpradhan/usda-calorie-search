# USDA Calorie Search

Go app that lists calorie information about foods, given a query.

## API

```
/api/v1/foods
```

Parameters:
- q: food query

Example Response:

```
URL: /api/v1/foods?q=raw%20onion

{
  "foods": Array[3][
    {
      "ndbno": "11282",
      "name": "Onions, raw",
      "weight": 160,
      "measure": "1.0 cup, chopped",
      "nutrients": Array[1][
        {
          "nutrient_id": "208",
          "nutrient": "Energy",
          "unit": "kcal",
          "value": "64",
          "gm": 40
        }
      ]
    },
    {
      "ndbno": "11291",
      "name": "Onions, spring or scallions (includes tops and bulb), raw",
      "weight": 100,
      "measure": "1.0 cup, chopped",
      "nutrients": Array[1][
        {
          "nutrient_id": "208",
          "nutrient": "Energy",
          "unit": "kcal",
          "value": "32",
          "gm": 32
        }
      ]
    },
    {
      "ndbno": "11294",
      "name": "Onions, sweet, raw",
      "weight": 148,
      "measure": "1.0 NLEA serving",
      "nutrients": Array[1][
        {
          "nutrient_id": "208",
          "nutrient": "Energy",
          "unit": "kcal",
          "value": "47",
          "gm": 32
        }
      ]
    }
  ]
}
```

