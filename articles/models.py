from django.db import models
from django.utils import timezone
from django.contrib.auth.models import User
from articles import fields

# Create your models here.
class Author(models.Model):
  name = models.CharField(max_length=100)
  created = models.DateTimeField('date registered', default = timezone.now)
  updated = fields.AutoDateTimeField('date updated', default = timezone.now)
  user = models.OneToOneField(User, related_name='author', unique=True)
  def __str__(self):
    return self.name

class Article(models.Model):
  title = models.CharField(max_length=200)
  created = models.DateTimeField('date published', default = timezone.now)
  updated = fields.AutoDateTimeField('date updated', default = timezone.now)
  author = models.ForeignKey(Author, related_name='articles')
  content = models.TextField()
  def __str__(self):
    return self.title+"\n\n"+self.content

class Anthology(models.Model):
  name = models.CharField(max_length=200)
  created = models.DateTimeField('date created', default = timezone.now)
  updated = fields.AutoDateTimeField('date updated', default = timezone.now)
  author = models.ForeignKey(Author, related_name='anthologies')
  articles = models.ManyToManyField(Article, related_name='anthologies', blank=True)
  def __str__(self):
    return self.name
