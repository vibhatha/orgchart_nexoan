# Description of sample data

Presidencies:

2025/01/01- 2025/01/31 President Ranil Wickremesinghe
2025/02/01- 2025/02/28: President Anura Kumara

Ranil

    Organisation:

        2025/01/3:
        - Add Minister of Defence: Sri Lanka Army, Sri Lanka Navy
        - Add Minister of Health: Sri Lanka Medical Council, Department of Ayurveda
        - Add Minister of Finance: General Treasury, Department of National Budget

        2025/01/15
        - Move General Treasury from Minister of Finance to Minister of Defence
        - Terminate Department of Ayurveda in Minister of Health

    People:

        2025/01/6:
        - Minister of Health: George Washington

        2025/01/10:
        - Minister of Finance: Hamilton

        2025/01/25:
        - Minister of Health: Kanye West

Anura

    Organisation:

        2025/02/05:
        - Add Minister of Defence: Sri Lanka Army, Sri Lanka Navy, Sri Lanka Air Force
        - Add Minister of Health: Sri Lanka Medical Council, Department of Ayurveda
        - Add Minister of Finance and Economy: General Treasury, Department of National Budget, Department of Public Finance

        2025/02/18
        - Move Sri Lanka Medical Council to Minister of Defence
        - Terminate Department of National Budget

    People:

        2025/02/12:
        - Minister of Finance and Economy: Vibhatha Abeykoon

        2025/02/14:
        - Minister of Health: Sanjiva Weerawarana
        - Minister of Finance and Economy: Ranil Wickremesinghe

        2025/02/26:
        - Minister of Health: Kanye West

To add all sample data:

For Ranil Wickremesinghe:

```bash
./orgchart -data $(pwd)/data/sample_data/presidents/2025-01-01/ -init -type person
./orgchart -data $(pwd)/data/sample_data/presidents/2025-01-31/ -type person
./orgchart -data $(pwd)/data/sample_data/presidents/2025-02-01/ -type person
./orgchart -data $(pwd)/data/sample_data/presidents/2025-02-28/ -type person

./orgchart -data $(pwd)/data/sample_data/rw/orgchart/2025-01-03 
./orgchart -data $(pwd)/data/sample_data/rw/orgchart/2025-01-15

./orgchart -data $(pwd)/data/sample_data/rw/people/2025-01-06 -type person
./orgchart -data $(pwd)/data/sample_data/rw/people/2025-01-10 -type person
./orgchart -data $(pwd)/data/sample_data/rw/people/2025-01-25 -type person
```

For Anura Kumara:

* Adding a new presidency does not work yet
 


