#version 330 core
out vec4 FragColor;
in vec2 TexCoord;

uniform float r;

void main(){
    FragColor = vec4(r,0.0f,0.5f,1.0f);
}
