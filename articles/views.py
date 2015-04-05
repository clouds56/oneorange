from django.shortcuts import render, redirect, get_object_or_404
from django.http import HttpResponse
from django.template import RequestContext, loader
from django.contrib import auth
from django.contrib.auth.decorators import login_required
from django.views.generic import TemplateView

from articles.models import Author, Article, Anthology
# Create your views here.

def index(request):
    authors = Author.objects.all()
    return render(request, "base.html", {'authors': authors})

class AuthorView(TemplateView):
    template_name = "author.html"

    def get(self, request, author_name):
        context = {}
        context['author'] = get_object_or_404(Author, name=author_name)
        context['author_anthologies'] = context['author'].anthologies.all()
        return self.render_to_response(context)

    def post(self, request, author_name):
        context = {}
        context['author'] = get_object_or_404(Author, name=author_name)
        anthology_name = ""
        created = None
        if "anthology_name" in request.POST:
            anthology_name = request.POST["anthology_name"]
        if anthology_name:
            anthology, created = Anthology.objects.get_or_create(name=anthology_name, author=context['author'])
            if created:
                context['msg'] = anthology_name + " created"
            else:
                context['msg'] = anthology_name + " already exist"
        else:
            context['msg'] = "error name"
        context['author_anthologies'] = context['author'].anthologies.all()
        return self.render_to_response(context)

def anthology(request, author_name, anthology_name):
    anthology = get_object_or_404(Anthology, author__name=author_name, name=anthology_name)
    return render(request, "anthology.html", {'anthology': anthology})

def article(request, author_name, anthology_name, article_id):
    article = get_object_or_404(Article, id=article_id)
    return render(request, "article.html", {'article': article})
