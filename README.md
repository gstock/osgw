# OSGW
## Documentación
**Endpoint**: `/api/osgw/{{username}}`

*username*: usuario de Github

**Response**

`{"count":30,"avg_temp":20.9}`

*count*: cantidad de repositorios

*avg_temp*: promedio de temperatura (solo de aquellos repositorios de los cuales se pudo obtener la temperatura)

## Notas

- Se decidió implementar como una **API GET** para facilitar las pruebas, pero podría haberse implementado como un POST con form-data o json body.
La API valida que el usuario no sea vacío. Sin embargo, dada la configuración de rutas esto no debería pasar. En este caso se devuelve: **400 BAD REQUEST**

- La validación de errores tiene varios IF anidados. Dado que Go no tiene exceptions, no pensé ni investigué en detalle las prácticas recomendadas.

- Si se produce un error al ejecutar la API se retorna un error de HTTP y en el body un json con la descripción del mensaje.
En el caso de que falle la conexión a GitHub, se devuelve: **503 SERVICE UNAVAILABLE**
Se puede pensar un poco más que códigos HTTP devolver en cada caso.
`{"error":"Error obtaining user info from Github","success":false}`

- En el caso de que falle la conexión a la API del clima, no se devuelve error. Esto se debe a que si un usuario tiene 100 repositorios y falla una o algunas llamadas a la api del clima, de todas formas se puede devolver un promedio de temperaturas. Se podría agregar en la respuesta un detalle donde se indiquen los errores.

- El campo location de Github es de texto libre. Dado que la API de World Weather Online puede recibir varios inputs (US Zipcode, UK Postcode, Canada Postalcode, IP address, Latitude/Longitude (decimal degree) or city name)), se decidió enviar este valor directamente. Sin embargo puede ser un problema, una mejora podría ser validar la location contra otro servicio.

- Para generar la URL utilice concatenación, ej.: httpclient.Get("https://api.github.com/users/" + username) , lo hice por practicidad, pero lo correcto sería generar un tipo URL donde se pueda ir generando la ruta o agregando los query params (No tuve tiempo de investigarlo en Go)

## Diseño

Es posible no haya utilizado correctamente las distintas estructuras de Go para organizar el código. La forma en que se pensó es la siguiente:

**App**: es un struct que contiene 2 interfaces a los clientes de GitHub y WorldWeatherOnline (WWO). Se pensó de esta forma para poder armar Mocks en los test cases.
El método `getRepoAvgTemp` de App es quien recibe la request, verifica los parametros y obtiene el usuario y los repositorios desde GitHub. Luego recorre todos los repositorios y obtiene la temperatura para cada fecha de creación. Termino siendo un método largo, y lo ideal sería extraer toda esta lógica a otro método o incluso clase.

**GitClient**: Cliente que se conecta a Github. Implementa la interfaz RepoClient con 2 métodos: 
```
GetUser(string) (User, error)
GetRepos(string) ([]Repository, error)
```
GitClient devuelve dos modelos:

- **User**: representa un usuario de GitHub
- **Repository**: representa un repositorio

*Solo se implementaron los campos necesarios

**WWOClient**: Cliente que se conecta a WWO. Implementa la interfaz WeatherClient con el método:
```
GetWeather(string, string) (Weather, error)
```

## Mejoras

Dado que las temperaturas historicas no cambian, se puede guardar en una base de datos los datos "Location", "Date", "Temperature" de forma tal de no utilizar la API de clima en todos los casos. Antes de hacer esto sería conveniente unificar las location, es decir que: "Buenos Aires", "buenos aires" y "Buenos Aires, Argentina" mapeen a una misma Location de nuestra base local, de esta forma se ahorrarían más llamadas a la API.

Tambien se podría cachear por un periodo de tiempo los datos del usuario de github (location y repositorios), pero estos datos pueden cambiar.

## Tests

`api_test.go` contiene los siguientes tests:

- TestGetRepoAvgTempNoUser: usuario vacío
- TestGetRepoAvgTemp: usuario de Github existente
- TestInvalidUser: usuario de Github no existente (aplica a falló de conexión con GitHub)
- TestValidUser: usuario válido y con repos
- TestNoRepos: usuario válido y sin repos (no enconte información de que devuelve GitHub si no hay repos, supuse un array vacío)
- TestValidUserInvalidDate: usuario válido, con repos, con una fecha invalida.

Existen Mocks para GitClient y WWOCLient

## Swagger
Inclui una pequeña descripcion de la API con Swagger, corriendo `genswagger.sh` se genera la documentación e inicia el servidor para visualizarla en el browser. Esta es la razon por la que se habilito CORS.
