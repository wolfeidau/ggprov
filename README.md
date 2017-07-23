# ggprov 

The ggprov project is a suite of tools to provision, configure and deploy [AWS greengrass](https://aws.amazon.com/greengrass/) to a linaro based platform such as the [Dragonboard 410c](https://developer.qualcomm.com/hardware/dragonboard-410c).

# overview

The project consists of a few command line tools which are:

* gg-prov which provisions all the greengrass and IoT related resources, then exports a configuration file.
* gg-deploy which copies the deployable files to the linaro based system.
* gg-config which installs and configures greengrass on the linaro based system.

# disclaimer

This is a work in progress at the moment.

# licence 

This code is released under MIT License, and is copyright Mark Wolfe.
