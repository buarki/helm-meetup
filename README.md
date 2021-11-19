# Helm

This is a pratical quick straightforward hands-on tutorial about Helm. You can and must look for deeper theory and topics after this one.

## TOC

- [1. What is helm](#1-what-is-helm)
- [2. Creating a very simple chart](#2-creating-a-very-simple-chart)
- [3. How to tear down a release](#3-how-to-tear-down-a-release)
- [4. Linting and testing k8s objects](#4-linting-and-testing-k8s-objects)
- [5. Add comments](#5-add-commands)
- [6. Using multiple Values.yaml file](#6-using-multiple-values-file)
- [7. Using templates](#7-using-templates)
    * [7.1 Getting started](#71-getting-started)
    * [7.2 Conditionals](#72-conditionals)
    * [7.3 More about conditionals and useful functions](#73-more-about-conditionals-and-useful-functions)
- [8. Subcharts](#8-subcharts)
    * [8.1 A chart under charts directory](#81-a-chart-under-charts-directory)
    * [8.2 Global values for main chart and children](#82-global-values-for-main-chart-and-children)
    * [8.3 Importing external charts](#83-importing-external-charts)
    * [8.4 Overriding child chart values](#84-overriding-child-chart-values)
- [9. Packaging and publishing charts](#9-packaging-and-publishing-charts)
    * [9.1 Setting up Chart Museum](#91-setting-up-chart-museum)
    * [9.2 List available repos](#92-list-available-registry)
    * [9.3 Package a chart](#93-package-a-chart)
    * [9.4 Publish a packaged chart](#94-publishing-a-packaged-chart)
    * [9.5 Searching charts](#95-searching-charts)
    * [9.6 Instaling from registry](#96-instaling-from-registry)
    * [9.7 Instaling specific version](#97-instaling-specific-version)
    * [9.8 Deleting a chart from registry](#98-deleting-a-chart-from-registry)
    * [9.9 Deleting a chart registry](#99-deleting-a-chart-registry)
  

## 1. What is helm <div id='1-what-is-helm'>


[Helm](https://helm.sh/) is a tool that helps us manage k8s applications by doing versioning, publishing and it also enables sharing. All of that is done by charts. 

Helm charts are just yaml files with necessary setup to deploy an application on a kubernetes cluster.

## 2. Creating a very simple chart <div id='2-creating-a-very-simple-chart'>

For this tutorial I suggest you to use minikube to make things easy. So before proceeding with this one, install and start minikube.

To create a chart called `app-1` run the following command:

```
helm create app-1
```

Once you ran that command you'll see a directory with name `app-1`. If you open that directory you'll see folders `charts` and `templates` and three files called `.helmignore`, `Chart.yaml` and `values.yaml`. 

*FOR NOW, JUST REMOVE DIRECTORY `charts` AND UNDER TEMPLATE DIRECTORY REMOVE ALL FILES EXCEPT `deployment.yaml`, `ingress.yaml` and `service.yaml`. AND ALSO, DELETE ALL `CONTENT` OF FILE `values.yaml` (NOT THE FILE)*.

Once you remove those unnecessary files, replace the content of `deployment.yaml` with the following:

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Values.deployment.name }}
spec:
  replicas: {{ .Values.deployment.replicas }}
  selector:
    matchLabels:
      serviceName: {{ .Values.deployment.serviceName }}
  template:
    metadata:
      labels:
        serviceName: {{ .Values.deployment.serviceName }}
    spec:
      containers:
        - name: app-1
          image: "{{ .Values.deployment.image.repository }}:{{ .Values.deployment.image.tag }}"
          ports:
            - containerPort: {{ .Values.port }}
```

Then, add the following content to `Values.yaml`:

```
# this is the port that container will use
port: 8080

deployment:
  name: app-1
  replicas: 1
  serviceName: app-1
  image:
    repository: oreasek/app-1
    tag: latest
```

That syntax of `"{{ }}"` comes from [Go templates](https://blog.gopheracademy.com/advent-2017/using-go-templates/). Summarizing what it does, it is just used to mark places that we want to replace by somthing. In the `deployment.yaml` file, for example, we are saying that the `replicas` spec attribute will be provided by the `replicas` of `deployment` object defined at `Values.yaml` file. The same way, `containerPort` on `deployment.yaml` will be replaced by `containerPort` defined at `deployment` object at `Values.yaml` file.

The idea is that: on object manifest we mark something that we want to provide by a "values file" and we provide on a "values" file.

Extending that ideia, replace the content of `service.yaml` file by the following:

```
apiVersion: v1
kind: Service
metadata:
  name: {{ .Values.service.name }}
spec:
  type: {{ .Values.service.type }}
  selector:
    serviceName: {{ .Values.service.serviceName }}
  ports:
    - port: {{ .Values.port }}
      targetPort: {{ .Values.port }}
```

Then *APPEND* the following content to `Values.yaml`:

```
service:
  name: app-1-service
  type: ClusterIP
  serviceName: app-1
```

Also replace the content of `ingress.yaml` by the following:

```
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ .Values.ingress.name }}
  annotations:
    kubernetes.io/ingress.class: nginx
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
    - host: "{{ .Values.ingress.host }}"
      http:
        paths:
          - path: /
            pathType: ImplementationSpecific
            backend:
              service:
                name: {{ .Values.service.name }}
                port: 
                  number: {{ .Values.port }}
```

And also *APPEND* the following content to `Values.yaml`:

```
ingress:
  name: app-1-ingress
  host: app-1-using-helm.com
```

In the end, `Values.yaml` will look like this:

```
# this is the port that container will use
port: 8080

deployment:
  name: app-1
  replicas: 1
  serviceName: app-1
  image:
    repository: oreasek/app-1
    tag: latest

service:
  name: app-1-service
  type: ClusterIP
  serviceName: app-1

ingress:
  name: app-1-ingress
  host: app-1-using-helm.com
```

Before we proceed, also add the host `app-1-using-helm.com` pointing to minikube ip to `/etc/hosts`:

```
sudo echo "$(minikube ip) app-1-using-helm.com" >> /etc/hosts
```

Once you have changed the content of `deployment.yaml`, `service.yaml`, `ingress.yaml` and `Values.yaml` run the following command:

```
helm install app-1 .
```

If you didn't get any error right now you should have a service running. To check it using helm you can do:

```
helm ls
```

And you'll see something like that:

```
NAME 	NAMESPACE	REVISION	UPDATED                                	STATUS  	CHART      	APP VERSION
app-1	default  	1       	2021-11-14 18:13:44.137298776 -0300 -03	deployed	app-1-0.1.0	1.16.0  
```

If you want the output in JSON form you can do:

```
helm ls -o json
```

```
[{"name":"app-1","namespace":"default","revision":"1","updated":"2021-11-14 18:13:44.137298776 -0300 -03","status":"deployed","chart":"app-1-0.1.0","app_version":"1.16.0"}]
```

To check if the service is really running you can do:

```
curl https://app-1-using-helm.com -k
```

And you should see as output `"palmeiras nao tem mundial"`.


You can also check if the Deployment configuration matches with what was requested:

```
kubectl describe deployment app-1
```

The result will look like this:

```
Name:                   app-1
Namespace:              default
CreationTimestamp:      Sun, 14 Nov 2021 18:13:44 -0300
Labels:                 app.kubernetes.io/managed-by=Helm
Annotations:            deployment.kubernetes.io/revision: 1
                        meta.helm.sh/release-name: app-1
                        meta.helm.sh/release-namespace: default
Selector:               serviceName=app-1
Replicas:               1 desired | 1 updated | 1 total | 1 available | 0 unavailable
StrategyType:           RollingUpdate
MinReadySeconds:        0
RollingUpdateStrategy:  25% max unavailable, 25% max surge
Pod Template:
  Labels:  serviceName=app-1
  Containers:
   app-1:
    Image:        oreasek/app-1:latest
    Port:         8080/TCP
    Host Port:    0/TCP
    Environment:  <none>
    Mounts:       <none>
  Volumes:        <none>
Conditions:
  Type           Status  Reason
  ----           ------  ------
  Available      True    MinimumReplicasAvailable
  Progressing    True    NewReplicaSetAvailable
OldReplicaSets:  <none>
NewReplicaSet:   app-1-7977d4cd7c (1/1 replicas created)
Events:          <none>
```

That output shows us that it should have one replica, use port 8080 and run image `oreasek/app-1:latest`.


## 3. How to tear down a release <div id='3-how-to-tear-down-a-release'>

If you need to fully remove a released service, you can just run:

```
helm uninstall app-1
```

## 4. Linting and testing k8s objects <div id='4-linting-and-testing-k8s-objects'>

To check if *the yaml object is valid in terms of yaml only* you can run:

```
helm lint
```

To check if generated k8s object and inspect them you can run:

```
helm template .
```

To check if generated k8s object and test if they are correct using kubectl you can run:

```
helm template . | kubectl apply -f - --dry-run=client
```

## 5. Add comments <div id='5-add-commands'>
If you want to comment something on `Values.yaml` or on any other k8s manifests file you can do it by using `#`.

If you want to add a comment at file `_helpers.tpl`, which is used to add pure template stuff, not related at all with k8s, you can use that syntax:

```
{{/*

a comment here

*/}}
```

## 6. Using multiple Values.yaml file <div id='6-using-multiple-values-file'>

As shown above, the `Values.yaml` file is used to define values to be injected on the manifests files. It is possible to create one file like that per environment, e.g: one for production environment and another one for testing environment.

To check it working, create another file called `test.yaml` and add the following content:

```
# this is the port that container will use
port: 8080

deployment:
  name: app-1
  replicas: 5
  serviceName: app-1
  image:
    repository: oreasek/app-1
    tag: test

service:
  name: app-1-service
  type: ClusterIP
  serviceName: app-1

ingress:
  name: app-1-ingress
  host: app-1-using-helm-test.com
```

NOTE: Remember to add host `app-1-using-helm-test.com` to your /etc/hosts file.

The difference on this new file is that it says that deployment for app-1 must have 5 replicas and should use tag `test`.

To deploy that service using that configs you can run:

```
helm install app-1-test . -f test.yaml
```

Once it gets up, if you run:

```
curl https://app-1-using-helm-test.com -k
```

You should get exact same result as before.

## 7. Using templates <div id='7-using-templates'>

### 7.1 Getting started <div id='71-getting-started'>

Create a file called `_app1templates.tpl` under `templates` directory.

We're gonna do a refactor on our chart already using templates. If you get back to `deployment.yaml` you'll see that we have labels 
`serviceName`. Those labels are used on deployment definition but also on service definition at file `service.yaml`. So we are somehow duplicating code. So lets refact the chart in such a way we define those labels only in one place.

On the created `_app1templates.tpl` file add the following:

```
{{- define "app-1-labels" -}}
serviceName: app-1
{{- end -}}
```

With above code we are creating a way to reference our application labels.


Replace the content of `service.yaml` by the following:

```
apiVersion: v1
kind: Service
metadata:
  name: {{ .Values.service.name }}
  labels: {{- include "app-1-labels" . | nindent 4 }}
spec:
  type: {{ .Values.service.type }}
  selector: {{- include "app-1-labels" . | nindent 4 }}
  ports:
    - port: {{ .Values.port }}
      targetPort: {{ .Values.port }}
```


*KEEP KALM! YOU'LL UNDERSTAND WHAT THOSE NEW LINES DO.* First of all, lets inspect the following line: 

```
selector: {{- include "app-1-labels" . | nindent 4 }}
```

What it does is:
- the first `-` after the `{{` means: break this line;
- the "`include "app-1-labels" .`" means: grab the content under `app-1-labels`, defined at file `_app1templates.tpl`, and place here;
- the "`| nindent 4`" means: move that content 4 spaces towards to left. This is necessary because if you look for the previous version of `service.yaml` you'll see that it had 4 spaces of distance between the label defined there and the left border.


We also added `labels` attribute at deployment's metadata:

```
labels: {{- include "app-1-labels" . | nindent 4 }}
```

And it has the same effect of previous one.

Replace the content of `deployment.yaml` by bellow one:

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Values.deployment.name }}
  labels: {{- include "app-1-labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.deployment.replicas }}
  selector:
    matchLabels: {{- include "app-1-labels" . | nindent 6 }}
  template:
    metadata:
      labels: {{- include "app-1-labels" . | nindent 8 }}
    spec:
      containers:
        - name: app-1
          image: "{{ .Values.deployment.image.repository }}:{{ .Values.deployment.image.tag }}"
          ports:
            - containerPort: {{ .Values.port }}
```

The explanation for those updates are the same of what we did with `services.yaml`, so take a time and understand they.

To see the how that chart looks like when we "compile" it, run the following command:

```
helm template .
```

Output will look like this:

```
---
# Source: app-1/templates/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: app-1-service
  labels:
    serviceName: app-1
spec:
  type: ClusterIP
  selector:
    serviceName: app-1
  ports:
    - port: 8080
      targetPort: 8080
---
# Source: app-1/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-1
  labels:
    serviceName: app-1
spec:
  replicas: 1
  selector:
    matchLabels:
      serviceName: app-1
  template:
    metadata:
      labels:
        serviceName: app-1
    spec:
      containers:
        - name: app-1
          image: "oreasek/app-1:latest"
          ports:
            - containerPort: 8080
---
# Source: app-1/templates/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app-1-ingress
  annotations:
    kubernetes.io/ingress.class: nginx
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
    - host: "app-1-using-helm.com"
      http:
        paths:
          - path: /
            pathType: ImplementationSpecific
            backend:
              service:
                name: app-1-service
                port: 
                  number: 8080
```

### 7.2 Conditionals <div id='72-conditionals'>

We are able to build k8s manifests conditionally based in some configurations. The image `oreasek/app-1` accept an environment variable called `MESSAGE`, if that one is not provided a default one will be used. To see how conditionals work, *APPEND* to `Values.yaml` file the following content:

```
containerEnv:
  override: true
```

Then, update the content of `deployment.yaml` file to bellow one:

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Values.deployment.name }}
  labels: {{- include "app-1-labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.deployment.replicas }}
  selector:
    matchLabels: {{- include "app-1-labels" . | nindent 6 }}
  template:
    metadata:
      labels: {{- include "app-1-labels" . | nindent 8 }}
    spec:
      containers:
        - name: app-1
          image: "{{ .Values.deployment.image.repository }}:{{ .Values.deployment.image.tag }}"
          ports:
            - containerPort: {{ .Values.port }}
          {{- if eq .Values.containerEnv.override true }}
          env:
            - name: MESSAGE
              value: https://www.youtube.com/watch?v=IfZZcHubtto
          {{- end -}}
```

What was done:
- we added a condition to `deployment.yaml` that checks if variable `override` of object `containerEnv` defined at `Values.yaml` is true or not;
- if it is true, `env` attribute will be overriden and curl should return `https://www.youtube.com/watch?v=IfZZcHubtto`;
- if it is false nothing changes


Uinstall `app-1` chart if it is installed:

```
helm uninstall app-1
```

Then, deploy the chart:

```
helm install app-1 .
```

If you run:

```
curl https://app-1-using-helm.com -k
```

You should see "https://www.youtube.com/watch?v=IfZZcHubtto"

### 7.3 More about conditionals and useful functions <div id='73-more-about-conditionals-and-useful-functions'>

More info about conditionals: https://helm.sh/docs/chart_template_guide/control_structures/;

About useful functions: https://helm.sh/docs/chart_template_guide/function_list/#network-functions;

## 8. Subcharts <div id='8-subcharts'>

Sometimes is useful to group related charts together. This can be done by placing charts under `charts` directory.

### 8.1 A chart under charts direct <div id='81-a-chart-under-charts-directory'>

On `app-1`'s directory, recreate a directory called `charts`

```
mkdir charts
```

Under the new created charts directory, create a new chart:

```
helm create mlclient
```

The `mlclient` will be a service that will expose an endpoint and provide answears based on request. It will be client of a "machine learning" service.

On "templates" directory delete all files except `deployment.yaml`, `ingress.yaml` and `service.yaml`.

Replace the content of `deployment.yaml` with:

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Values.app.name }}
spec:
  replicas: {{ .Values.deployment.replicas }}
  selector:
    matchLabels:
      app: {{ .Values.app.name }}
  template:
    metadata:
      labels:
        app: {{ .Values.app.name }}
    spec:
      containers:
        - name: mlclient
          image: "{{ .Values.deployment.image.repository }}:{{ .Values.deployment.image.tag }}"
          ports:
            - containerPort: {{ .Values.app.port }}
          env:
            - name: ML_SERVICE_HOST
              value: {{ .Values.deployment.env.ML_SERVICE_HOST }}
```

Replace the content of `service.yaml` with:

```
apiVersion: v1
kind: Service
metadata:
  name: {{ .Values.app.name }}
spec:
  type: {{ .Values.service.type }}
  selector:
    app: {{ .Values.app.name }}
  ports:
    - port: {{ .Values.app.port }}
      targetPort: {{ .Values.app.port }}
```

Replace the content of `ingress.yaml` with:

```
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ .Values.app.name }}
  annotations:
    kubernetes.io/ingress.class: nginx
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
    - host: "{{ .Values.ingress.host }}"
      http:
        paths:
          - path: /
            pathType: ImplementationSpecific
            backend:
              service:
                name: {{ .Values.app.name }}
                port: 
                  number: {{ .Values.app.port }}
```

And then, replace the content of `values.yaml` with:

```
app:
  name: mlclient
  port: 8080

deployment:
  replicas: 1
  image:
    repository: oreasek/mlclient
    tag: latest

service:
  type: ClusterIP
  
ingress:
  host: mlclient.com
```

Also, add `mlclient.com` to your /etc/hosts:

```
echo $(minikube up) mlclient.com >> /etc/hosts
```

Once you did it, deploy with helm:

```
helm install mlclient .
```

Check installation by:

```
curl https://mlclient.com -k
```

If output is "ask something" everything is fine. For now uninstall the chart by:

```
helm uninstall mlclient .
```

As said before, `mlclient` is as client of `mlservice`, so we need to setup it. Let's setup it as *subchart*. To do it, under directory `charts` run:

```
helm create mlservice
```

On that new created chart remove all files from `templates` directory except from `deployment.yaml`, `service.yaml` and `ingress.yaml`.

Once you removed those unnecessary files, replace the content of `deployment.yaml` to:

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Values.app.name }}
spec:
  replicas: {{ .Values.deployment.replicas }}
  selector:
    matchLabels:
      app: {{ .Values.app.name }}
  template:
    metadata:
      labels:
        app: {{ .Values.app.name }}
    spec:
      containers:
        - name: mlservice
          image: "{{ .Values.deployment.image.repository }}:{{ .Values.deployment.image.tag }}"
          ports:
            - containerPort: {{ .Values.app.port }}
```

And replace the content of `service.yaml` with:

```
apiVersion: v1
kind: Service
metadata:
  name: {{ .Values.app.name }}
spec:
  type: {{ .Values.service.type }}
  selector:
    app: {{ .Values.app.name }}
  ports:
    - port: {{ .Values.app.port }}
      targetPort: {{ .Values.app.port }}
```

Then, replace the content of `ingress.yaml` with:

```
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ .Values.app.name }}
  annotations:
    kubernetes.io/ingress.class: nginx
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
    - host: "{{ .Values.ingress.host }}"
      http:
        paths:
          - path: /
            pathType: ImplementationSpecific
            backend:
              service:
                name: {{ .Values.app.name }}
                port: 
                  number: {{ .Values.app.port }}
```

And the content of `values.yaml` with:

```
app:
  name: mlservice
  port: 8081

deployment:
  replicas: 1
  image:
    repository: oreasek/mlservice
    tag: latest

service:
  type: ClusterIP
  
ingress:
  host: mlservice.com
```

Also, add host `mlservice.com` to your `/etc/hosts`:

```
echo $(minikube ip) mlservice.com >> /etc/hosts
```

Once you did all that you can deploy the `mlclient` with its subchart `mlservice`:


```
helm install mlclient .
```

Once it get deployed you ask a question, for example:

```
curl -G https://mlclient.com -k --data-urlencode "question=is Pele better than Maradona"
```

### 8.2 Global values for main chart and children <div id='82-global-values-for-main-chart-and-children'>

Sometimes the main chart and its children have values almost the same, in such a way it is good define those values only once. 

We can do it using global values defined at `values.yaml` *OF MAIN CHART*.

If you compare k8s manifests of `mlclient` and `mlservice` you see that they have same values for:
 - service type (ClusterIP);
 - deployment replicas (1);
 - image tag (latest);

Due to that, we can group them, so replace the content of `values.yaml` file *OF MAIN CHART* with this content:

```
global:
  deployment:
    replicas: 1
    image:
      tag: latest
  service:
    type: ClusterIP

app:
  name: mlclient
  port: 8080

deployment:
  image:
    repository: oreasek/mlclient
  env:
    ML_SERVICE_HOST: http://mlservice.default.svc.cluster.local:8081
  
ingress:
  host: mlclient.com
```

Then, replace the content of *MAIN CHART* `deployment.yaml` with:

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Values.app.name }}
spec:
  replicas: {{ .Values.global.deployment.replicas }}
  selector:
    matchLabels:
      app: {{ .Values.app.name }}
  template:
    metadata:
      labels:
        app: {{ .Values.app.name }}
    spec:
      containers:
        - name: mlclient
          image: "{{ .Values.deployment.image.repository }}:{{ .Values.global.deployment.image.tag }}"
          ports:
            - containerPort: {{ .Values.app.port }}
          env:
            - name: ML_SERVICE_HOST
              value: {{ .Values.deployment.env.ML_SERVICE_HOST }}
```

Replace the content of *SUBCHART* `deployment.yaml` with:

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Values.app.name }}
spec:
  replicas: {{ .Values.global.deployment.replicas }}
  selector:
    matchLabels:
      app: {{ .Values.app.name }}
  template:
    metadata:
      labels:
        app: {{ .Values.app.name }}
    spec:
      containers:
        - name: mlservice
          image: "{{ .Values.deployment.image.repository }}:{{ .Values.global.deployment.image.tag }}"
          ports:
            - containerPort: {{ .Values.app.port }}
```

Replace the content of *MAIN CHART* `service.yaml` with:

```
apiVersion: v1
kind: Service
metadata:
  name: {{ .Values.app.name }}
spec:
  type: {{ .Values.global.service.type }}
  selector:
    app: {{ .Values.app.name }}
  ports:
    - port: {{ .Values.app.port }}
      targetPort: {{ .Values.app.port }}
```

Replace the content of *SUBCHART* `service.yaml` with:

```
apiVersion: v1
kind: Service
metadata:
  name: {{ .Values.app.name }}
spec:
  type: {{ .Values.global.service.type }}
  selector:
    app: {{ .Values.app.name }}
  ports:
    - port: {{ .Values.app.port }}
      targetPort: {{ .Values.app.port }}
```

Once you did it, uninstall the chart and deploy it again:

```
helm install mlclient .
```

Then, check it by runnig a curl:

```
curl -G https://mlclient.com -k --data-urlencode "question=is go better than kotling?"
```

If you got an answear things are fine.

You can checkout this example at directory *subcharts-present*.

### 8.3 Importing external charts <div id='83-importing-external-charts'>

Another way to have subcharts is importing charts, as well as importing Go modules.

To check it create a new chart on a different directory than previous example.

Then, create a new chart:

```
helm create amainchart
```

On the root of created folder `amainchart` create a new file called `requirements.yaml`.

Then add a the repository of bitnami:

```
helm repo add bitnami https://charts.bitnami.com/bitnami
```

Then, add bellow code to `requirements.yaml`:

```
dependencies:
  - name: redis
    version: 15.5.5
    repository: https://charts.bitnami.com/bitnami
```

That yaml code is saying that:
- we want to use chart called `redis`;
- with chart version 15.5.5;
- from repository https://charts.bitnami.com/bitnami;

You can also list dependencies of a chart by running:

```
helm dep ls
```

To "install" the external chart you should run:

```
helm dep build
```

You'll see that a tar file will be added at `charts` directory and also a `Chart.lock` file will be created. The `Chart.lock` works like a `yarn.lock` (it freezes a charts version).

To deploy it you can run:

```
helm install myredis .
```

### 8.4 Overriding child chart values <div id='84-overriding-child-chart-values'>

The default chart `bitnami/redis` has some interesting configs like a master node, but in our case we want the simple version with only one node. 

Firs of all, uninstall previous chart:

```
helm uninstall amainchart
```

To override the redis chart you just need to add an object with same name of chart, ´redis´ in this case to `values.yaml`. You can see how the file looks like [here](importing-external-chart/amainchart/values.yaml).

Sumarizing what that file does, it is overriding the chart to use redis in stand alone mode.


## 9. Packaging and publishing charts <div id='9-packaging-and-publishing-charts'>

### 9.1 Setting up Chart Museum <div id='91-setting-up-chart-museum'>

[Chart Museum](https://github.com/helm/chartmuseum) allows you to store and serve Helm Charts.

To run it using helm, first of all, add the repo for Chart Museum:

```
helm repo add chartmuseum https://chartmuseum.github.io/charts
```

Then, update repositories:

```
helm repo update
```

On a different directory of the one you was using to follow this tutorial, run the following command:

```
helm pull chartmuseum/chartmuseum
```

Once helm pull that chart, the chart is pulled as a compressed file, so unzip it:

```
tar -xvf $NAME_OF_DOWNLOADED_CHART
```

Then, go to the created directory. If everything went fine now you can see the chart files. We need to adjust the `Values.yaml`, so:

- look for `DISABLE_API` and make it false;
- change the `service` definition to its type be a `NodePort` instead of a `ClusterIP` to make the setup easy. And also define the `nodePort` for a high value, like 31515.

Then, install the Chart Museum:

```
helm install chartmuseum .
```

If everything went fine you can open your browser and search for `$MINIKUBEIP:$THE_PORT_YOU_DEFINED_FOR_NODE_PORT`. If you see "Welcome to ChartMuseum!" everything is fine.

Once Chart Museum is up, you should add it as a repo. The following command adds it with name `localregistry`:

```
helm repo add localregistry "http://$(minikube ip):$(kubectl get --namespace default -o jsonpath="{.spec.ports[0].nodePort}" services chartmuseum)"
```

### 9.2 List available repos <div id='92-list-available-registry'>

Check if registry `localregistry` is present on available resgistries list.

```
helm repo list
```

If it is, update repos:

```
helm repo update
```

### 9.3 Package a chart <div id='93-package-a-chart'>

Go to the directory where you was running the tutorial and run:

```
helm package .
```

Now you should see a zipped file with the name of your chart.

### 9.4 Publish a packaged chart <div id='94-publishing-a-packaged-chart'>

To publish a chart to registry you can do it using curl:

```
curl --data-binary "@app-1-0.1.0.tgz" "http://$(minikube ip):$(kubectl get --namespace default -o jsonpath="{.spec.ports[0].nodePort}" services chartmuseum)/api/charts"
```

If you get message "{"saved":true}" the chart has been published.

### 9.5 Searching charts <div id='95-searching-charts'>

First, update repos:

```
helm repo update
```

To search for a chart, run:

```
helm search repo app-1
```

### 9.6 Instaling from registry <div id='96-instaling-from-registry'>

Once you know the chart name, you can install it direct from registry. In this case using the published `app-1` to install it with a new name of `app-1-from-registry` just run:

```
helm install app-1-from-registry localregistry/app-1
```

If installation went fine you can check by running a curl:

```
curl https://app-1-using-helm.com -k
```

### 9.7 Instaling specific version <div id='97-instaling-specific-version'>


If you want to specify version to install run (uninstall previous before):
helm install app-1-from-registry localregistry/app-1 --version 

```
helm install app-1-from-registry-with-version localregistry/app-1 --version 0.1.0
```

### 9.8 Deleting a chart from registry <div id='98-deleting-a-chart-from-registry'>

To delete a chart you need to do DELETE specifying chart and version to delete:

```
curl -X DELETE "http://$(minikube ip):$(kubectl get --namespace default -o jsonpath="{.spec.ports[0].nodePort}" services chartmuseum)/api/charts/app-1/0.1.0"
```

### 9.9 Deleting a chart registry <div id='99-deleting-a-chart-registry'>

To delete a chart registry just run:

```
helm repo remove localregistry
```
