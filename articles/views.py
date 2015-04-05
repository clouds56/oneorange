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

def author(request, author_name):
    args = {}
    args['author'] = get_object_or_404(Author, name=author_name)
    if request.method == "GET":
        args['author_anthologies'] = args['author'].anthologies.all()
        return render(request, "author.html", args)
    elif request.method == "POST":
        anthology_name = ""
        created = None
        if "name" in request.POST:
            anthology_name = request.POST["name"]
        if anthology_name:
            anthology, created = Anthology.objects.get_or_create(name=anthology_name, author=args['author'])
        else:
            args['author_anthologies'] = args['author'].anthologies.all()
            args['msg'] = "error name"
            render(request, "author.html", args)
        if created:
            args['msg'] = anthology_name + " created"
        else:
            args['msg'] = anthology_name + " already exist"
        args['author_anthologies'] = args['author'].anthologies.all()
        return render(request, "author.html", args)

def anthology(request, author_name, anthology_name):
    anthology = get_object_or_404(Anthology, author__name=author_name, name=anthology_name)
    return render(request, "anthology.html", {'anthology': anthology})

def article(request, author_name, anthology_name, article_id):
    article = get_object_or_404(Article, id=article_id)
    return render(request, "article.html", {'article': article})
