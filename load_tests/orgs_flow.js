import http from 'k6/http';
import { check } from 'k6';

// Конфигурация теста
export let options = {
  vus: 333, // 999 RPS
  duration: '1m',
};

function randomEmail() {
  const timestamp = Date.now(); // время в миллисекундах
  return `testuser_${timestamp}_${Math.floor(Math.random() * 100000)}@mail.com`;
}


function randomString(length) {
  const chars = 'abcdefghijklmnopqrstuvwxyz';
  let result = '';
  for (let i = 0; i < length; ++i) {
    result += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return result;
}

export default function () {
  // 1. Регистрация
  let orgReq = {
    uuid: "",
    email: randomEmail(),
    password: "SuperSecurePassword1!",
    name: "Test Organization",
    rating: 4.5,
    address: "123 Test Street",
    type: "clinic",
    telephone: "+79990001122",
    city: "Moscow",
    about: "Test org description",
    lat: 55.7558,
    long: 37.6173,
  };

  let orgPayload = {
    UUID: orgReq.uuid,
    email: orgReq.email,
    password: orgReq.password,
    name: orgReq.name,
    rating: orgReq.rating,
    address: orgReq.address,
    type: orgReq.type,
    telephone: orgReq.telephone,
    city: orgReq.city,
    about: orgReq.about,
    lat: orgReq.lat,
    long: orgReq.long,
  };

  let res = http.post('http://localhost:8100/v1/auth/registration/orgs', JSON.stringify(orgPayload), {
    headers: { 'Content-Type': 'application/json' },
  });

  check(res, {
    'registration status is 200': (r) => r.status === 200,
  });

  let tokens;
  try {
    tokens = JSON.parse(res.body);
  } catch (e) {
    console.error('Ошибка разбора токенов: ', e);
    return;
  }

  let accessToken = tokens.access_token;
  if (!accessToken) return;

  let authHeaders = {
    headers: {
      Authorization: `Bearer ${accessToken}`,
      'Content-Type': 'application/json',
    },
  };

  // 2. Добавление работника
  let workerPayload = {
    worker_info: {
      first_name: "Иван",
      last_name: "Тестов",
      position: "Врач",
    },
  };

  let workerRes = http.post('http://localhost:8100/v1/orgs/workers', JSON.stringify(workerPayload), authHeaders);
  check(workerRes, {
    'worker creation status is 200': (r) => r.status === 200,
  });

  // 3. Добавление услуги
  let servicePayload = {
    service_info: {
      name: "Прием эндокринолога",
      cost: 2000,
      description: "Первичный осмотр врача",
    },
  };

  let serviceRes = http.post('http://localhost:8100/v1/orgs/services', JSON.stringify(servicePayload), authHeaders);
  check(serviceRes, {
    'service creation status is 200': (r) => r.status === 200,
  });
}
