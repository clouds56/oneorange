from django.shortcuts import render
from django.http import HttpResponse
from django.template import RequestContext, loader
from django.contrib import auth

from articles.models import Author, Article, Anthology
# Create your views here.

def index(request):
    authors = Author.objects.all()
    template = loader.get_template("base.html")
    context = RequestContext(request, {'authors': authors})
    return HttpResponse(template.render(context))

def login(request):
    username = request.POST['username']
    password = request.POST['password']
    user = auth.authenticate(username=username, password=password)
    if user is not None:
        auth.login(request, user)
        return HttpResponse('/accounts/loggedin')
    else:
        return HttpResponse('/accounts/invalid')

def author(request, author):
    anthologies = Anthology.objects.filter(author__name=author);
    template = loader.get_template("author.html")
    context = RequestContext(request, {'anthologies': anthologies})
    return HttpResponse(template.render(context))

def anthology(request, author, anthology):
    articles = Article.objects.filter(anthologies__name=anthology);
    template = loader.get_template("anthology.html")
    context = RequestContext(request, {'articles': articles, 'anthology_name': anthology})
    return HttpResponse(template.render(context))

def article(request, author, anthology, article):
    article = Article.objects.get(id=article)
    template = loader.get_template("article.html")
    context = RequestContext(request, {'article': article})
    return HttpResponse(template.render(context))
