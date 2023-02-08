# go-raytrace 

This repository implements a ray tracer in Go by following [The Ray Tracer Challenge](https://pragprog.com/titles/jbtracer/the-ray-tracer-challenge/) book
* The book lays out a language-agnostic TDD plan for implementing a ray tracer from scratch
* Window rendering and input handling provided by Go library [ebiten](https://github.com/hajimehoshi/ebiten)
* I highly recommend the book. 5 stars!

## Controls
```
W: Zoom in
S: Zoom out
A/D/Q/E/Z/C: Move camera position (wip)
LeftArrow/UpArrow/RightArrow/DownArrow: Move camera focal point (wip)
N: next scene
M: next camera position
Numpad plus: Increase ray bounces
Numpad minus: Decrease ray bounces
Numpad multiply (*): Increase rendering goroutines
Numpad divide (/): Decrease rendering goroutines
```

---

## Examples
*Bonus online chapters: Soft shadows, Bounding Volume Hierarchies, UV texture mapping, Skyboxes*
![Example 1](renders/tori_marbles_skybox.png?raw=true "bonus chapter example 1")
![Example 2](renders/groups_skybox.png?raw=true "bonus chapter example 2")
http://www.raytracerchallenge.com/bonus/area-light.html  
http://raytracerchallenge.com/bonus/bounding-boxes.html  
http://www.raytracerchallenge.com/bonus/texture-mapping.html  

*Chapter 15: Triangles & Meshes & Parse .obj files. Bonus feature: parse .rpl Toribash replay files*  
![Chapter 15 example](renders/file_parsing.png?raw=true "chapter 15 example")

Scene pictures one of my [Toribash fight](https://www.youtube.com/watch?v=zpGLPsczHGU&t=46s) replays immortalized into seven glass marbles. With T-posing mesh reptilians in a staring contest.  
Toribash replays only capture keyframes with the engine simulating the rest in between, some frames are lost in the render.  
Send me a replay file if you want a scene of your own! 
- Add .obj file parsing. Creates raytracing triangle natives from mesh description.
- Add .rpl Toribash replay file parsing. Creates raytracing natives from replay keyframes.


*Chapter 14: Groups*

![Chapter 14 example](renders/group_transforms.png?raw=true "chapter 14 example")

*Chapter 13: Cylinders and cones*

![Chapter 13 example](renders/cone_and_cylinder.png?raw=true "chapter 13 example")

*Chapter 11: Spheres and mirror*

![Chapter 11 example](renders/spheres_mirror.png?raw=true "chapter 11 example")
  
*Chapter 11: Hollow glass sphere against checkerboard*

![Chapter 11 example](renders/hollow_glass_sphere.png?raw=true "chapter 11 example")

*Chapter 11: Submerged spheres in water*

![Chapter 11 example](renders/water.png?raw=true "chapter 11 example")

*Chapter 4: Rendering shapes*

![Chapter 4 example](renders/shapes.png?raw=true "chapter 4 example")


