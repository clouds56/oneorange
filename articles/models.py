from django.db import models

# Create your models here.
class Author(models.Model):
  name = models.CharField(max_length=100)
  reg_date = models.DateTimeField('date registered')
  def __str__(self):
    return self.name

class Article(models.Model):
  title = models.CharField(max_length=200)
  content = models.CharField(max_length=20000)
  pub_date = models.DateTimeField('date published')
  author = models.ForeignKey(Author)
  def __str__(self):
    return self.title+"\n\n"+self.content
