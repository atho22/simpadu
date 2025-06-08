# Tutorial Penggunaan API Simpadu di Laravel

## Persiapan Awal

1. Buat project Laravel baru:
```bash
composer create-project laravel/laravel simpadu-client
cd simpadu-client
```

2. Install package yang diperlukan:
```bash
composer require guzzlehttp/guzzle
```

3. Buat file konfigurasi untuk API Simpadu di `config/simpadu.php`:
```php
<?php

return [
    'base_url' => env('SIMPADU_API_URL', 'http://localhost:8080'),
    'timeout' => 30,
];
```

4. Tambahkan konfigurasi di `.env`:
```env
SIMPADU_API_URL=http://localhost:8080
```

## Implementasi Service

1. Buat service class untuk menangani API calls di `app/Services/SimpaduService.php`:
```php
<?php

namespace App\Services;

use GuzzleHttp\Client;
use Illuminate\Support\Facades\Session;

class SimpaduService
{
    protected $client;
    protected $baseUrl;

    public function __construct()
    {
        $this->baseUrl = config('simpadu.base_url');
        $this->client = new Client([
            'base_uri' => $this->baseUrl,
            'timeout' => config('simpadu.timeout'),
        ]);
    }

    protected function getHeaders()
    {
        $headers = [
            'Accept' => 'application/json',
            'Content-Type' => 'application/json',
        ];

        // Tambahkan token jika ada
        if (Session::has('simpadu_token')) {
            $headers['Authorization'] = 'Bearer ' . Session::get('simpadu_token');
        }

        return $headers;
    }

    public function login($username, $password)
    {
        try {
            $response = $this->client->post('/login', [
                'headers' => $this->getHeaders(),
                'json' => [
                    'username' => $username,
                    'password' => $password,
                ],
            ]);

            $data = json_decode($response->getBody(), true);
            
            // Simpan token
            if (isset($data['token'])) {
                Session::put('simpadu_token', $data['token']);
            }

            return $data;
        } catch (\Exception $e) {
            throw new \Exception('Login gagal: ' . $e->getMessage());
        }
    }

    public function getProfile()
    {
        try {
            $response = $this->client->get('/profile', [
                'headers' => $this->getHeaders(),
            ]);

            return json_decode($response->getBody(), true);
        } catch (\Exception $e) {
            throw new \Exception('Gagal mengambil profil: ' . $e->getMessage());
        }
    }

    // Tambahkan method lain sesuai kebutuhan
}
```

## Implementasi Controller

1. Buat AuthController di `app/Http/Controllers/AuthController.php`:
```php
<?php

namespace App\Http\Controllers;

use App\Services\SimpaduService;
use Illuminate\Http\Request;

class AuthController extends Controller
{
    protected $simpaduService;

    public function __construct(SimpaduService $simpaduService)
    {
        $this->simpaduService = $simpaduService;
    }

    public function showLogin()
    {
        return view('auth.login');
    }

    public function login(Request $request)
    {
        $request->validate([
            'username' => 'required',
            'password' => 'required',
        ]);

        try {
            $response = $this->simpaduService->login(
                $request->username,
                $request->password
            );

            return redirect()->route('dashboard')->with('success', 'Login berhasil');
        } catch (\Exception $e) {
            return back()->with('error', $e->getMessage());
        }
    }

    public function logout()
    {
        session()->forget('simpadu_token');
        return redirect()->route('login');
    }
}
```

## Implementasi View

1. Buat view login di `resources/views/auth/login.blade.php`:
```html
<!DOCTYPE html>
<html>
<head>
    <title>Login Simpadu</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <div class="container mt-5">
        <div class="row justify-content-center">
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header">Login Simpadu</div>
                    <div class="card-body">
                        @if(session('error'))
                            <div class="alert alert-danger">
                                {{ session('error') }}
                            </div>
                        @endif

                        <form method="POST" action="{{ route('login') }}">
                            @csrf
                            <div class="mb-3">
                                <label>Username</label>
                                <input type="text" name="username" class="form-control" required>
                            </div>
                            <div class="mb-3">
                                <label>Password</label>
                                <input type="password" name="password" class="form-control" required>
                            </div>
                            <button type="submit" class="btn btn-primary">Login</button>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
</body>
</html>
```

## Implementasi Routes

1. Tambahkan routes di `routes/web.php`:
```php
<?php

use App\Http\Controllers\AuthController;
use Illuminate\Support\Facades\Route;

Route::get('/', function () {
    return redirect()->route('login');
});

Route::get('/login', [AuthController::class, 'showLogin'])->name('login');
Route::post('/login', [AuthController::class, 'login']);
Route::post('/logout', [AuthController::class, 'logout'])->name('logout');

// Protected routes
Route::middleware(['auth.simpadu'])->group(function () {
    Route::get('/dashboard', function () {
        return view('dashboard');
    })->name('dashboard');
});
```

## Implementasi Middleware

1. Buat middleware untuk autentikasi di `app/Http/Middleware/SimpaduAuth.php`:
```php
<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Session;

class SimpaduAuth
{
    public function handle(Request $request, Closure $next)
    {
        if (!Session::has('simpadu_token')) {
            return redirect()->route('login');
        }

        return $next($request);
    }
}
```

2. Daftarkan middleware di `app/Http/Kernel.php`:
```php
protected $routeMiddleware = [
    // ...
    'auth.simpadu' => \App\Http\Middleware\SimpaduAuth::class,
];
```

## Implementasi JavaScript untuk Local Storage

1. Buat file JavaScript di `public/js/auth.js`:
```javascript
// Fungsi untuk menyimpan token
function saveToken(token) {
    localStorage.setItem('simpadu_token', token);
}

// Fungsi untuk mengambil token
function getToken() {
    return localStorage.getItem('simpadu_token');
}

// Fungsi untuk menghapus token
function removeToken() {
    localStorage.removeItem('simpadu_token');
}

// Fungsi untuk mengecek apakah user sudah login
function isLoggedIn() {
    return getToken() !== null;
}

// Fungsi untuk logout
function logout() {
    removeToken();
    window.location.href = '/login';
}

// Fungsi untuk menambahkan token ke header
function addTokenToHeader(headers = {}) {
    const token = getToken();
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }
    return headers;
}

// Contoh penggunaan dengan fetch API
async function fetchWithAuth(url, options = {}) {
    const headers = addTokenToHeader(options.headers || {});
    
    try {
        const response = await fetch(url, {
            ...options,
            headers: {
                ...headers,
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            }
        });

        if (response.status === 401) {
            // Token expired atau tidak valid
            removeToken();
            window.location.href = '/login';
            return;
        }

        return await response.json();
    } catch (error) {
        console.error('Error:', error);
        throw error;
    }
}

// Contoh penggunaan untuk login
async function login(username, password) {
    try {
        const response = await fetch('/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify({ username, password })
        });

        const data = await response.json();
        
        if (data.token) {
            saveToken(data.token);
            window.location.href = '/dashboard';
        } else {
            throw new Error('Token tidak ditemukan dalam response');
        }
    } catch (error) {
        console.error('Login error:', error);
        throw error;
    }
}
```

2. Tambahkan script ke view login:
```html
<script src="{{ asset('js/auth.js') }}"></script>
```

## Contoh Penggunaan API

1. Mengambil data jadwal:
```php
// Di SimpaduService.php
public function getJadwalKelas($kelasId)
{
    try {
        $response = $this->client->get("/jadwal/kelas/{$kelasId}", [
            'headers' => $this->getHeaders(),
        ]);

        return json_decode($response->getBody(), true);
    } catch (\Exception $e) {
        throw new \Exception('Gagal mengambil jadwal: ' . $e->getMessage());
    }
}

// Di Controller
public function jadwalKelas($kelasId)
{
    try {
        $jadwal = $this->simpaduService->getJadwalKelas($kelasId);
        return view('jadwal.index', ['jadwal' => $jadwal['data']]);
    } catch (\Exception $e) {
        return back()->with('error', $e->getMessage());
    }
}
```

2. Mengambil data ruangan:
```php
// Di SimpaduService.php
public function getRuanganByGedung($gedung)
{
    try {
        $response = $this->client->get("/ruangan/gedung/{$gedung}", [
            'headers' => $this->getHeaders(),
        ]);

        return json_decode($response->getBody(), true);
    } catch (\Exception $e) {
        throw new \Exception('Gagal mengambil data ruangan: ' . $e->getMessage());
    }
}

// Di Controller
public function ruanganGedung($gedung)
{
    try {
        $ruangan = $this->simpaduService->getRuanganByGedung($gedung);
        return view('ruangan.index', ['ruangan' => $ruangan['data']]);
    } catch (\Exception $e) {
        return back()->with('error', $e->getMessage());
    }
}
```

## Catatan Penting

1. **Keamanan**:
   - Jangan simpan token di localStorage untuk aplikasi yang memerlukan keamanan tinggi
   - Gunakan HttpOnly cookies untuk menyimpan token
   - Implementasikan refresh token untuk keamanan yang lebih baik

2. **Error Handling**:
   - Selalu tangani error dengan baik
   - Berikan feedback yang jelas ke user
   - Log error untuk debugging

3. **Best Practices**:
   - Gunakan environment variables untuk konfigurasi
   - Implementasikan rate limiting
   - Cache response yang sering digunakan
   - Validasi input sebelum mengirim ke API

4. **Testing**:
   - Buat unit test untuk service
   - Buat feature test untuk controller
   - Mock API calls saat testing

## Langkah Selanjutnya

1. Implementasikan refresh token
2. Tambahkan validasi input yang lebih lengkap
3. Implementasikan error handling yang lebih baik
4. Tambahkan fitur "Remember Me"
5. Implementasikan rate limiting
6. Tambahkan logging untuk debugging
7. Implementasikan caching untuk response yang sering digunakan 