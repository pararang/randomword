


1. Clone dan buat branch dgn nama masing2
2. Buat folder dgn nama masing2
3. di dalam folder buat aplikasi/service dengan satu endpoint yang mengembalikan hasil dari service yg ada dibawah
    - host: https://random-indonesian-word.p.rapidapi.com [GET]
    - path: /words/random
    - parameter: 
      - limit[int] (optional)
    - header
      - X-RapidAPI-Host: random-indonesian-word.p.rapidapi.com
      - X-RapidAPI-Key: e9b079b792msh3e1aa8608df520dp15d272jsn0831149d10d6
4. spesifikasi endpoint:
  - GET
  - /words
  - Response: JSON
    ```json
    {"words": ["kata1", "kata2"]}
    ```
