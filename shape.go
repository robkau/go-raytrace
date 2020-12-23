package main

type shape interface {
	intersect(ray) intersections
	normalAt(tuple) tuple
	getTransform() x4Matrix
	setTransform(matrix x4Matrix) shape
	getMaterial() material
	setMaterial(material) shape
}

func normalAt(p tuple, t x4Matrix, lnaf func(tuple) tuple) tuple {
	localPoint := t.invert().mulTuple(p)
	localNormal := lnaf(localPoint)
	worldNormal := t.invert().transpose().mulTuple(localNormal)
	worldNormal.c = vector
	return worldNormal.normalize()
}

func intersect(r ray, t x4Matrix, lif func(ray) intersections) intersections {
	lr := r.transform(t.invert())
	return lif(lr)
}
