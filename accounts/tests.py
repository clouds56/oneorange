from django.test import TestCase, Client
from django.contrib.auth.models import User
import re

# Create your tests here.

class AccountsTest(TestCase):

    def setUp(self):
        user = User.objects.create_user('zz', 'zz@example.com', '123')
        user.save()
#    def test_for_test(self):
#        self.assertFalse(User.objects.all().exists())

    def test_login(self):
        client = Client()
        response = client.get('/accounts/login')
        self.assertRedirects(response, '/accounts/login?next=/articles', status_code=302)

        response = client.get('/accounts/login/')
        self.assertRedirects(response, '/accounts/login?next=/articles', status_code=302)

        response = client.post('/accounts/login', {'password': '123'})
        self.assertContains(response, '<p>No username</p>', 1)
        self.assertContains(response, '<input type="hidden" name="next" value="/articles">', 1)


        response = client.post('/accounts/login', {'username': 'zz', 'next': '/abc'})
        #print(re.findall('.*hidden.*', response.content.decode(), re.MULTILINE))
        self.assertContains(response, '<p>No password</p>', 1)
        self.assertContains(response, '<input type="hidden" name="next" value="/abc">', 1)

        response = client.post('/accounts/login/?next=/abc', {'username': 'zz', 'password': '12'})
        self.assertContains(response, '<p>Wrong username or password</p>', 1)
        self.assertContains(response, '<input type="hidden" name="next" value="/abc">', 1)

        response = client.post('/accounts/login?', {'username': 'zz', 'password': '123', 'next': '/accounts/signup'})
        self.assertRedirects(response, '/accounts/signup', status_code=302)

    def test_logout(self):
        client = Client()
        response = client.get('/accounts/logout')
        self.assertRedirects(response, '/accounts/login/?next=/accounts/logout', status_code=302)

        response = client.post('/accounts/login?', {'username': 'zz', 'password': '123', 'next': '/accounts/login/?next=/accounts/logout'})
        response = client.get('/accounts/logout', follow=True)
        self.assertTrue('articles' in str(response.redirect_chain))