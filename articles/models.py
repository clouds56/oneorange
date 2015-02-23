from django.db import models
from django.utils import timezone
from django.contrib.auth.models import User

# Create your models here.
class Author(models.Model):
  name = models.CharField(max_length=100)
  user = models.OneToOneField(User, related_name='author', unique=True)
  date_create = models.DateTimeField('date registered', default = timezone.now)
  def __str__(self):
    return self.name

class Article(models.Model):
  title = models.CharField(max_length=200)
  author = models.ForeignKey(Author, related_name='articles')
  content = models.TextField()
  date_create = models.DateTimeField('date published', default = timezone.now)
  def __str__(self):
    return self.title+"\n\n"+self.content

class Anthology(models.Model):
  name = models.CharField(max_length=200)
  date_create = models.DateTimeField('date created', default = timezone.now)
  author = models.ForeignKey(Author, related_name='anthologies')
  articles = models.ManyToManyField(Article, related_name='anthologies')
  def __str__(self):
    return self.name
