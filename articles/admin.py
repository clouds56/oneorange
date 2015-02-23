from django.contrib import admin
from articles.models import Author, Article, Anthology
# Register your models here.

class ArticleAdmin(admin.ModelAdmin):
    fieldsets = [
        (None, {'fields': ['title']}),
        ('Info', {'fields': ['author', 'date_create'], 'classes': ['collapse']}),
        (None, {'fields': ['content']}),
    ]
    list_display = ('title', 'author', 'date_create')

admin.site.register(Author)
admin.site.register(Anthology)
admin.site.register(Article, ArticleAdmin)
