<h1>Lazy Persistence</h1>
Esta é a versão do projeto Lazy Persistence escrita na linguagem Go. A reescrita do projeto visa melhorar a performance do componente,
aumentando o throughput e diminuindo o consumo de memória da máquina.
Diferente da sua versão em Java, a nova versão do Lazy Persistence nunca bloquea, em alternativa ao modelo de sincronização de threads, 
foram utilizados conceitos do paradigma funcional criando cópias dos componentes que necessitam de acesso concorrente. 


<h2>Funcionalidades Adicionais</h2>
<ul>
<li>Suporte à múltiplas filas</li>
<li>Suporte à múltiplos databases</li>
<li>Envio de metadados no Header da mensagem</li>
<li>Nova action increment</li>
<li>Suporte à operações idempotentes</li>
</ul>

<h2>Instalação</h2>
O projeto utiliza o GoVendor como ferramenta de gerenciamento de dependência, o mesmo pode ser obtído com o comando abaixo. 

<pre>
go get -u github.com/kardianos/govendor
</pre>

Execute o seguinte comando após clonar o repositório do projeto

<pre>
govendor sync
</pre>

O comando irá baixar as dependências do projeto. 

<h2>Exemplo de Configuração</h2>

```yaml
rabbitmq: 
  host: localhost 
  port: 5672
  user: guest
  password: guest
  queues:
    - queue: 
      name: lazypersistence
      exchange: monitriip 
      exchangetype: topic 
      routingkey: monitriip.logs
      dlqexchange: lazypersistence.DLQ
      dlqexchangetype: direct 
      dlqroutingkey: erros
      durable: true 
    - queue: 
      name: transmissoes
      durable: true
    - queue: 
      name: metricas
      exchange: frota 
      exchangetype: direct 
      durable: true 

mongodb: 
  host: localhost 
  port: 27017
  database: test

maxmessages: 300
retentionseconds: 30
logfile: /home/rodrigo/lazypersistence.log
```
O arquivo de configuração interno fica dentro da pasta <b>m2m-lazypersistence/configs</b>

<h2>Ambientes</h2>
O projeto suporta os seguintes ambientes de execução: 

<ul>
<li>PRODUCTION</li>
<li>DEVELOPMENT (default)</li>
</ul>

O ambiente de execução pode ser configurado através da variável de ambiente M2M-ENVIRONMENT ou através de uma flag em tempo de execução.

<h2>Execução</h2>

Compilação do projeto
<pre>
go build -o lazypersistence cmd/lazypersistence/main.go
</pre>
Execução
<pre>
./lazypersistence -config-location=/home/dev/config.yml -m2m-environment=PRODUCTION
</pre>
Executando em desenvolvimento 
<pre>
go run cmd/lazypersistence/main.go
</pre>

<h2>Utilizando o serviço</h2>
A forma de utilizar o novo Lazypersistence é bastante similar a versão anterior, o serviço continua utilizando o protocolo AMQP porém os metadados não poluem
mais o payload da mensagem, os mesmos são enviados pelo header. 

O novo Lazypersistence suporta os seguintes metadados: <br/><br/>

<table>
<tr>
<th>Nome</th><th>Descrição</th><th>Exemplo</th>
</tr>
<tr>
<td>action</td><td>Define a ação a ser realizada sobre a mensagem enviada.</td><td>insert,push,pull ou increment</td>
</tr>
<tr>
<td>database</td><td>Define o bando de dados onde o documento será persistido.</td><td>frota_znd, monitriip_znd, etc</td>
</tr>
<tr>
<td>collection</td><td>Define a coleção onde o documento será persistido.</td><td>Transmissao, Bilhetes, etc</td>
<tr/>
<tr>
<td>field</td><td>Pode ser utilizado em conjunto com as actions push e pull e define o campo do documento onde essas operações serão realizadas. 
O campo deverá ser um array</td><td>viagens</td>
</tr>
<tr>
<td>id</td><td>Permite selecionar um documento pelo id, definindo o documento que será afetado pelas alterações, o metadado pode ser utilizam em conjunto com as actions push, pull e increment.</td>
<td>5adf8d3b1ddeb5b732f5caf5</td>
</tr>
<tr>
<td>condition</td><td>Permite selecionar um documento com base em um critério de busca. O metadado pode ser utilizado em conjunto com as actions push, pull e increment.</td>
<td>{"nome":"jose","idade":14}</td>
</tr>
<tr>
<td>create</td><td>Define se o documento será persistido caso o mesmo não exista. O metadado pode ser utilizado em conjunto com as actions push, pull, increment e o metadado de seleção condition.</td>
<td>true ou false</td>
</tr>
</table>

<h2>Ações</h2>

- Increment <br/>
A action increment serve para incrementar valores de forma preguiçosa em um documento, os valores de incremento são passados no corpo da mensagem. 
Considere a seguinte mensagem AMQP: 

<b>Headers</b><br/>
action: increment<br/>
condition: {"empresa":"m2m"}<br/>
collection: empresas<br/>
database: frota_znd<br/>
create: true <br/>
<b>Payload</b><br/>
{"qtdProjetos":1,"qtdFuncionarios":5,"qtdDepartamentos":2}<br/>

O Lazypersistence irá incrementar os valores passados nas colunas referenciadas no payload da mensagem e criar o documento caso o mesmo ainda não exista.
O documento criada no mongo será algo parecido com este: 

<pre>
{"qtdProjetos":1,"qtdFuncionarios":5,"qtdDepartamentos":2,"empresa":"m2m"}
</pre>
