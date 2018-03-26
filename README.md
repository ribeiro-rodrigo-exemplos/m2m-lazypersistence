<h1>Lazy Persistence</h1>
Esta é a versão do projeto Lazy Persistence escrita na linguagem Go. A reescrita do projeto visa melhorar a performance do componente,
aumentando o throughput e diminuindo o consumo de memória da máquina.
Diferente da sua versão em Java, a nova versão do Lazy Persistence nunca bloquea, em alternativa ao modelo de sincronização de threads, 
foram utilizados conceitos do paradigma funcional criando cópias dos componentes que necessitam de acesso concorrente. 


<h2>Funcionalidades Adicionais</h2>
<ul>
<li>Suporte à múltiplas filas</li>
<li>Suporte à múltiplos databases</li>
</ul>

<h2>Instalação</h2>
O projeto utiliza o GoVendor como ferramenta de gerenciamento de dependência, o mesmo pode ser obtído com o comando abaixo. 

<pre>
go get -u github.com/kardianos/govendor
</pre>

Execute os seguintes comandos após clonar o repositório do projeto

<pre>
govendor init
</pre>
<pre>
govendor fetch +missing
</pre>

Os comandos irão baixar as dependências do projeto. 

<h2>Configuração</h2>

```yaml
rabbitmq: 
  host: localhost 
  port: 5672
  user: guest
  password: guest
  queues:
    - lazypersistence
    - transmissoes

mongodb: 
  host: localhost 
  port: 27017
  database: test

maxmessages: 30
retentionseconds: 30
logfile: /home/dev/lazypersistence.log
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
