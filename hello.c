
#include <pthread.h>
#include <stdio.h> 

int i = 0;

void* thread_1Func() {
	int j;
	for (j = 0; j < 1000000; j++) {
		i++;
	}
	
	return NULL;
}

void* thread_2Func() {
	int k;
	for (k = 0; k < 1000000; k++) {
		i--;
	}
	return NULL;
}


int main(){
	
	pthread_t thread_1;
	pthread_t thread_2;
	
	pthread_create(&thread_1, NULL, thread_1Func, NULL);
	

	pthread_create(&thread_2, NULL, thread_2Func, NULL);
	pthread_join(thread_1, NULL);
	pthread_join(thread_2, NULL);

	printf("%d\n", i);

	return 0;
}
