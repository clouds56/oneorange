from django.db import models
from datetime import datetime

# Create your models here.
class Author(models.Model):
  name = models.CharField(primary_key=True, max_length=100)
  reg_date = models.DateTimeField('date registered')
  def __str__(self):
    return self.name

class Article(models.Model):
  title = models.CharField(max_length=200)
  author = models.ForeignKey(Author, related_name='articles')
  content = models.TextField()
  pub_date = models.DateTimeField('date published', default = datetime.now)
  def __str__(self):
    return self.title+"\n\n"+self.content

class Anthology(models.Model):
  name = models.CharField(max_length=200)
  create_date = models.DateTimeField('date created')
  author = models.ForeignKey(Author, related_name='anthologies')
  articles = models.ManyToManyField(Article, related_name='anthologies')
  def __str__(self):
    return self.name
