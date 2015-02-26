from django.shortcuts import render, redirect
from django.http import HttpResponse
from django.template import RequestContext, loader
from django.contrib import auth
from django.contrib.auth.models import User
from django.contrib.auth.decorators import login_required

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
    return redirect(request.META.get('HTTP_REFERER'))

@login_required
def logout(request):
    auth.logout(request)
    return redirect(request.META.get('HTTP_REFERER'))

def signup(request):
    if not 'csrfmiddlewaretoken' in request.POST:
        return render(request, "signup.html", {'error_msg': "New register"})
    if not 'username' in request.POST:
        return render(request, "signup.html", {'error_msg': "No username"})
    if not 'password' in request.POST:
        return render(request, "signup.html", {'error_msg': "No password"})
    if not 'email' in request.POST:
        return render(request, "signup.html", {'error_msg': "No email"})
    username = request.POST['username']
    email = request.POST['email']
    password = request.POST['password']
    if User.objects.filter(username=username).exists():
        return render(request, "signup.html", {'error_msg': "Username exists", 'last': {'username': username, 'email':email}})
    if User.objects.filter(email=email).exists():
        return render(request, "signup.html", {'error_msg': "Email exists", 'last': {'username': username, 'email':email}})
    user = User.objects.create_user(username, email, password)
    user.save()
    return redirect("/accounts/step2/"+username)

def step2(request, author):
    if not 'csrfmiddlewaretoken' in request.POST:
        return render(request, "step2.html", {'author': author})
    if request.user.is_authenticated():
        username = request.user.username
    if not request.POST['author']:
        return render(request, "step.html", {'error_msg': "No author name", 'author': author})
    author = request.POST['author']
    if Author.objects.filter(name=author).exists():
        return render(request, "step.html", {'error_msg': "Author name exists", 'author': author})
    anauthor = Author.objects.create(name=author, user=request.user)
    anauthor.save()
    return redirect("/accounts/success")
