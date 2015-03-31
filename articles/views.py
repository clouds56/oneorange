from django.shortcuts import render, redirect, get_object_or_404
from django.http import HttpResponse
from django.template import RequestContext, loader
from django.contrib import auth
from django.contrib.auth.decorators import login_required

from articles.models import Author, Article, Anthology
# Create your views here.

def index(request):
    authors = Author.objects.all()
    return render(request, "base.html", {'authors': authors})

def author(request, author):
    author = get_object_or_404(Author, name=author)
    author_anthologies = author.anthologies.all()
    return render(request, "author.html", {'author': author, 'author_anthologies': author_anthologies})

def anthology(request, author, anthology):
    anthology = get_object_or_404(Anthology, author__name=author, name=anthology)
    return render(request, "anthology.html", {'anthology': anthology})

def article(request, author, anthology, article):
    article = get_object_or_404(Article, id=article)
    return render(request, "article.html", {'article': article})
