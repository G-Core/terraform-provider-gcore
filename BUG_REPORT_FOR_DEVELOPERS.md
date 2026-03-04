# API Bug Report: PATCH /loadbalancers Returns Empty vrrp_ips

## English Version

### Summary
The `PATCH /cloud/v1/loadbalancers/{project_id}/{region_id}/{load_balancer_id}` endpoint returns incomplete data - specifically, `vrrp_ips` array is always empty `[]`, despite the OpenAPI specification indicating it should be populated.

### Impact
- Terraform provider fails with "element has vanished" error during LB updates
- SDK consumers cannot rely on PATCH response for computed fields
- Violates OpenAPI contract (spec says vrrp_ips should be included)

### Reproduction
```bash
# 1. Create LB
POST /cloud/v1/loadbalancers/{project_id}/{region_id}
Response: tasks[...] → wait → GET shows vrrp_ips: [2 elements] ✓

# 2. Rename LB
PATCH /cloud/v1/loadbalancers/{project_id}/{region_id}/{lb_id}
{"name": "new-name"}
Response: vrrp_ips: [] ❌ BUG!

# 3. GET LB
GET /cloud/v1/loadbalancers/{project_id}/{region_id}/{lb_id}
Response: vrrp_ips: [2 elements] ✓
```

**Demonstration script:** `demonstrate_vrrp_ips_issue.sh`

### Expected vs Actual

| Endpoint | vrrp_ips in response |
|----------|---------------------|
| GET | ✅ Populated (2 elements) |
| PATCH | ❌ Empty array `[]` |
| OpenAPI Spec | ✅ Says should be populated |

### Root Cause
PATCH endpoint likely returns the partial update object instead of re-fetching the complete resource from database.

### Recommended Fix

**Option 1 (BEST): Reuse GET logic**
```python
def patch(self, request, load_balancer_id):
    update_load_balancer(load_balancer_id, request.data)
    # Return the same response as GET would
    return self.get(request, load_balancer_id)
```

**Why:** Guarantees consistency between GET and PATCH responses, includes all computed fields, follows DRY principle.

### References
- OpenAPI Spec: `LoadbalancerSerializer` includes `vrrp_ips` (line 83296)
- PATCH definition: line 13231
- Jira: GCLOUD2-20778
- Reproduction: `/demonstrate_vrrp_ips_issue.sh`

---

## Русская версия

### Краткое описание
Endpoint `PATCH /cloud/v1/loadbalancers/{project_id}/{region_id}/{load_balancer_id}` возвращает неполные данные - конкретно, массив `vrrp_ips` всегда пустой `[]`, несмотря на то что OpenAPI спецификация указывает что он должен быть заполнен.

### Влияние
- Terraform provider падает с ошибкой "element has vanished" при обновлении LB
- Потребители SDK не могут полагаться на ответ PATCH для computed полей
- Нарушает OpenAPI контракт (спецификация говорит что vrrp_ips должен быть включен)

### Воспроизведение
```bash
# 1. Создать LB
POST /cloud/v1/loadbalancers/{project_id}/{region_id}
Ответ: tasks[...] → ждем → GET показывает vrrp_ips: [2 элемента] ✓

# 2. Переименовать LB
PATCH /cloud/v1/loadbalancers/{project_id}/{region_id}/{lb_id}
{"name": "new-name"}
Ответ: vrrp_ips: [] ❌ БАГ!

# 3. GET LB
GET /cloud/v1/loadbalancers/{project_id}/{region_id}/{lb_id}
Ответ: vrrp_ips: [2 элемента] ✓
```

**Скрипт демонстрации:** `demonstrate_vrrp_ips_issue.sh`

### Ожидаемое vs Фактическое

| Endpoint | vrrp_ips в ответе |
|----------|------------------|
| GET | ✅ Заполнен (2 элемента) |
| PATCH | ❌ Пустой массив `[]` |
| OpenAPI Spec | ✅ Указывает что должен быть заполнен |

### Первопричина
PATCH endpoint вероятно возвращает частичный update объект вместо того чтобы перечитать полный ресурс из базы данных.

### Рекомендуемое исправление

**Вариант 1 (ЛУЧШИЙ): Переиспользовать логику GET**
```python
def patch(self, request, load_balancer_id):
    update_load_balancer(load_balancer_id, request.data)
    # Вернуть такой же ответ как GET
    return self.get(request, load_balancer_id)
```

**Почему:** Гарантирует консистентность между GET и PATCH ответами, включает все computed поля, следует DRY принципу.

### Ссылки
- OpenAPI Spec: `LoadbalancerSerializer` включает `vrrp_ips` (строка 83296)
- PATCH определение: строка 13231
- Jira: GCLOUD2-20778
- Воспроизведение: `/demonstrate_vrrp_ips_issue.sh`

---

## Quick Copy-Paste for Slack/Email

**English:**
```
🐛 API Bug: PATCH /loadbalancers returns empty vrrp_ips[]

PATCH response: vrrp_ips: [] ❌
GET response: vrrp_ips: [2 elements] ✓

This breaks Terraform provider and violates OpenAPI spec.

Fix: Make PATCH reuse GET logic to return complete resource.
See: demonstrate_vrrp_ips_issue.sh for reproduction.

Jira: GCLOUD2-20778
```

**Russian:**
```
🐛 Баг API: PATCH /loadbalancers возвращает пустой vrrp_ips[]

PATCH ответ: vrrp_ips: [] ❌
GET ответ: vrrp_ips: [2 элемента] ✓

Это ломает Terraform provider и нарушает OpenAPI spec.

Исправление: PATCH должен переиспользовать логику GET для возврата полного ресурса.
См: demonstrate_vrrp_ips_issue.sh для воспроизведения.

Jira: GCLOUD2-20778
```
