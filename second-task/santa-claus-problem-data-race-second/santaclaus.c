#include <pthread.h>
#include <stdlib.h>
#include <assert.h>
#include <unistd.h>
#include <stdio.h>
#include <stdbool.h>

#include "semaphore.h"

pthread_t *CreateThread(void *(*f)(void *), void *a)
{
    pthread_t *t = malloc(sizeof(pthread_t));
    assert(t != NULL);
    int ret = pthread_create(t, NULL, f, a);
    assert(ret == 0);
    return t;
}

static const int N_ELVES = 10;      // Количество эльфов
static const int N_REINDEER = 9;    // Количество северных оленей

static int elves;                   // Текущее количество эльфов, ожидающих помощи
static int reindeer;                // Текущее количество северных оленей, ожидающих запряжки
static semaphore_t santaSem;        // Семафор для пробуждения Санты
static semaphore_t reindeerSem;     // Семафор для запряжки северных оленей
static semaphore_t elfTex;          // Семафор для ограничения количества эльфов, нуждающихся в помощи
static semaphore_t mutex;           // Мьютекс для защиты доступа к общим переменным

// Добавляем общую переменную для гонки данных
static int sharedData = 0;

void *SantaClaus(void *arg)
{
    printf("Santa Claus: Hoho, here I am\n");
    while (true)
    {
        Wait(santaSem); // Ожидание пробуждения Санты
        Wait(mutex); // Захват мьютекса
        if (reindeer == N_REINDEER)
        {
            printf("Santa Claus: preparing sleigh\n");
            for (int r = 0; r < N_REINDEER; r++)
                Release(reindeerSem); // Пробуждение всех оленей для запряжки
            printf("Santa Claus: make all kids in the world happy\n");
            reindeer = 0; // Сброс счетчика оленей
        }
        else if (elves == 3)
        {
            printf("Santa Claus: helping elves\n");
        }
        Release(mutex); // Освобождение мьютекса

        sharedData++; // Используем общую переменную без синхронизации
        printf("Santa Claus: shared data incremented to %d\n", sharedData);
    }
    return arg;
}

void *Reindeer(void *arg)
{
    int id = (int)arg;
    printf("This is reindeer %d\n", id);
    while (true)
    {
        Wait(mutex);
        reindeer++;
        if (reindeer == N_REINDEER)
            Release(santaSem); // Пробуждение Санты, если все олени собрались
        Release(mutex);
        Wait(reindeerSem); // Ожидание запряжки
        printf("Reindeer %d getting hitched\n", id);
        sleep(20); // Имитация времени запряжки

        // Используем общую переменную без синхронизации
        sharedData++;
        printf("Reindeer %d: shared data incremented to %d\n", id, sharedData);
    }
    return arg;
}

void *Elve(void *arg)
{
    int id = (int)arg;
    printf("This is elve %d\n", id);
    while (true)
    {
        bool need_help = random() % 100 < 10;
        if (need_help)
        {
            Wait(elfTex);
            Wait(mutex);
            elves++;
            if (elves == 3)
                Release(santaSem); // Пробуждение Санты, если 3 эльфа нуждаются в помощи
            else
                Release(elfTex);
            Release(mutex);

            printf("Elve %d will get help from Santa Claus\n", id);
            sleep(10); // Имитация времени помощи

            Release(mutex);
            elves--;
            if (elves == 0)
                Release(elfTex);
            Release(mutex);
        }

        printf("Elve %d at work\n", id);
        sleep(2 + random() % 5); // Имитация времени работы эльфа

        sharedData++; // Используем общую переменную без синхронизации
        printf("Elve %d: shared data incremented to %d\n", id, sharedData);
    }
    return arg;
}

int main(int ac, char **av)
{
    elves = 0;
    reindeer = 0;
    santaSem = CreateSemaphore(0);
    reindeerSem = CreateSemaphore(0);
    elfTex = CreateSemaphore(1);
    mutex = CreateSemaphore(1);

    pthread_t *santa_claus = CreateThread(SantaClaus, 0);

    pthread_t *reindeers[N_REINDEER];
    for (int r = 0; r < N_REINDEER; r++)
        reindeers[r] = CreateThread(Reindeer, (void *)r + 1);

    pthread_t *elves[N_ELVES];
    for (int e = 0; e < N_ELVES; e++)
        elves[e] = CreateThread(Elve, (void *)e + 1);

    int ret = pthread_join(*santa_claus, NULL);
    assert(ret == 0);
}