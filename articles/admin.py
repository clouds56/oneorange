from django.contrib import admin
from articles.models import Author, Article, Anthology
# Register your models here.

class ArticleAdmin(admin.ModelAdmin):
    fieldsets = [
        (None, {'fields': ['title']}),
        ('Info', {'fields': ['author', 'created', 'updated'], 'classes': ['collapse']}),
        (None, {'fields': ['content']}),
    ]
    list_display = ('title', 'author', 'created', 'updated')

admin.site.register(Author)
admin.site.register(Anthology)
admin.site.register(Article, ArticleAdmin)
