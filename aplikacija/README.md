# Pokretanje aplikacije

- Za pokretanje aplikacije potrebno je imati instaliran MongoDB (v6.0.12) i Go (v1.21.3)
- Postavljanje podataka za bazu je opisanu u radu
- Prije pokretanja potrebno je pokrenuti komandu `go get` da bi se pribavile vanjske biblioteke korištene u projektu
- Projekt se pokreće komandom `go run .`
- Web aplikacija dostupna na `localhost:8000` (u slučaju ako se nije mijenjala varijabla okoline PORT)

# Varijable okoline (optionalno)
- PORT (`8000`)

- MONGO (`default`):
    - MONGO_URI (`mongodb://0.0.0.0:27017/`)
    - MONGO_DB (`recommender`)
