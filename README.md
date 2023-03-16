# OpenAI GPT-3 Chatbot

This is a simple chatbot built using OpenAI's GPT-3. The chatbot can carry out a conversation with the user on a given topic using natural language processing. The chatbot will prompt you to pick an expert to chat with. Enter the name of the expert you wish to chat with, and then start chatting. Chatbot can remember the context of the chat as shown in usage example below. To exit the chat, type "stop" and press enter.


# Getting Started
To get started, you will need an OpenAI API key. You can sign up for one on the OpenAI website. Once you have your API key, add it to .env file.

Next, install the project dependencies by running:
```
go mod download
```

To run the chatbot, simply execute the following command:
```
go run main.go
```

# How It Works
The chatbot uses OpenAI's GPT-3 language model to generate responses to user input. When the user enters a message, the chatbot sends the conversation history to the GPT-3 API, which generates a response based on the context of the conversation so far. The chatbot then displays the response to the user and waits for the next input.

# Usage example

Can remember the context. This can be seen when "what is it good at ?"; context being golang. 

```
pick an expert: golang

golang expert here ask your question('stop' to end): what is golang ?
golang expert: Golang (short for Go) is an open source programming language developed by Google in 2007. It is designed to be simple, efficient, and easy to read and write. It is a statically typed language with garbage collection and provides built-in support for concurrency. Golang is often used for developing high-performance systems and distributed applications because of its fast compile times, low memory footprint, and excellent support for concurrent programming.


golang expert here ask your question ('stop' to end): what is it good at ?
golang expert: Golang is good at many things, some of which include:

1. Concurrent programming: Golang is designed with concurrency in mind, and it has built-in support for it. This makes it ideal for developing applications that need to handle many user requests at the same time.

2. Network programming: Golang has a powerful network programming package that allows developers to build high-performance network applications such as web servers, APIs, and microservices.

3. High-performance computing: Golang has a fast compiler and runtime, which makes it an excellent choice for developing high-performance computing applications, such as scientific computing, machine learning, and artificial intelligence.

4. Cloud computing: Golang has a rich set of libraries and frameworks that make it easy to build cloud-native applications, especially those that need to be deployed on platforms such as Kubernetes, Google Cloud Platform, and Amazon Web Services.

5. Cross-platform development: Golang allows developers to write code that can be compiled and run on multiple platforms, including Windows, Linux, macOS, and mobile devices. This makes it an attractive choice for developers who need to write cross-platform applications.


golang expert here ask your question('stop' to end):stop

```