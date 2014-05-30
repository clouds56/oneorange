from django.db import models
from datetime import datetime

# Create your models here.
class Author(models.Model):
  name = models.CharField(max_length=100)
  reg_date = models.DateTimeField('date registered')
  def __str__(self):
    return self.name

class Article(models.Model):
  title = models.CharField(max_length=200)
  author = models.ForeignKey(Author)
  content = models.TextField()
  pub_date = models.DateTimeField('date published', default = datetime.now)
  def __str__(self):
    return self.title+"\n\n"+self.content
