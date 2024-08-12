
# Table of Contents

1.  [Adding stuff](#org274f7ed)
2.  [ALU](#orgbe45918)



<a id="org274f7ed"></a>

# Adding stuff

Since we now have some Gates to work with, we can try to add things

![img](imgs/half-adder.png)

So we are continuing building Chips, but Chips that serve as mathematical functions, such as &ldquo;+&rdquo;.


<a id="orgbe45918"></a>

# ALU

Putting all of this together we arrive at our ALU, a Chip that is capable of calculating a multitude of functions. It takes as input two **16bit** numbers and has a **16bit** number as output (plus two Flags, one that is telling if the output was 0 (`zr`) and one that is telling if the output was negative (`ng`))

![img](imgs/ALU.png)

