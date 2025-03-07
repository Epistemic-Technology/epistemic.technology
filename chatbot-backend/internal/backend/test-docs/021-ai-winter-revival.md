# From AI Winter to Deep Learning Spring: 1990-2012

*Published: 2024-01-28*

The period known as the "AI Winter" marked a crucial turning point in artificial intelligence history. Following the collapse of the Japanese Fifth Generation Computer Project in 1992, which had invested ¥57 billion ($400 million) in AI research, funding for AI projects dropped by 74% worldwide between 1993 and 1994.

## The Winter Years

During the AI winter, even successful systems struggled. XCON, Digital Equipment Corporation's expert system for computer configuration, which had saved the company $40 million annually in the 1980s, required 308 maintenance programmers by 1990 to manage its 31,000 rules. The system's maintenance costs had ballooned to $5 million per year, exemplifying the scalability problems of classical AI approaches.

## Early Signs of Spring

The development of Support Vector Machines by Vladimir Vapnik in 1995 marked the beginning of the thaw. Running on a Sun SPARCstation 20, the first SVM implementation could process 50,000 training examples with 20 features in 12 minutes, achieving 98.9% accuracy on the MNIST dataset using only 2MB of RAM, a remarkable efficiency for the time.

## The Deep Learning Revolution

The breakthrough came in 2006 with Hinton's work on deep belief networks. Using a cluster of 10 Pentium 4 processors (2.4 GHz each), his team trained a network with 3 hidden layers of 500 units each, achieving 98.75% accuracy on MNIST after 259,200 CPU-hours of training. The network required only 1.6 million parameters, compared to the 60 million typically needed by previous approaches.

## Hardware Evolution

The revival was accelerated by hardware advances. The first use of GPUs for neural network training in 2007 using NVIDIA's G80 architecture demonstrated a 72x speedup over CPU implementations. A single GPU could now train a convolutional network in 4 days instead of 8 months, processing 9 million parameters using 128MB of GPU memory.

The transition from winter to spring was marked by both theoretical advances and practical implementations. By 2012, when AlexNet won the ImageNet competition with a 15.3% error rate (compared to 26.2% for the second-best entry), the field had completely transformed. AlexNet's training on two GTX 580 GPUs took 5 days to process 60 million parameters, a task that would have required 3 years on the fastest systems available during the AI winter. 