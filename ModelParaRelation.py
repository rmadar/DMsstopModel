import numpy    as np
from   math import pi,sqrt


######################################################################
#
#
#  Author: Romain Madar (romain.madar@cern.ch)
#  Date  : 27/07/17 
#
#  Dark matter model coupled to top quarks
#     arXiv:1407.7529 / JHEP01 (2015) 017
#
#  Non-resonnant model with following details/assumptions:
#   - largrangian density  eq(2.25)
#   - coupling between first and third generation only (FCNC)
#   - aL=0, aR(*)=gSM and gchiL(*)=gchiR(*)=gDM
#   - partial widths are expressed in eq (3.1) for (in)visible decay
#
#  Purpose of this code:
#   These functions allow to explore the model phenomenology: how the
#   total width and the BR into invisible change with couplings (and
#   vis-versa)?
#
#
######################################################################



# Top mass
#---------
mt = 172


# Phase space function for invisible decay
#-----------------------------------------
def PhiInv(mDM,mV):
    r=mDM/mV
    return  (mV/(12*pi))*  sqrt(1-4*r**2) * (1+2*r**2)

PhiInv = np.vectorize(PhiInv) # Enable numpy array for mV


# Phase space function for visible decay
#---------------------------------------
def PhiVis(mt,mV):
    r=mt/mV
    return (mV/pi) * (1-r**2) * (1-0.5*r**2-0.5*r**4)

PhiVis = np.vectorize(PhiVis) # Enable numpy array for mV


# Visible width
#--------------
def get_width_vis(gSM,mV):
    return gSM**2 * PhiVis(mt,mV)


# Invisible width
#----------------
def get_width_inv(gDM,mV,mDM):
    return gDM**2 * PhiInv(mDM,mV)


# Total width
#------------
def get_total_width(gSM, gDM, mV, mDM):
    return get_width_vis(gSM,mV) + get_width_inv(gDM,mV,mDM)


# Invisible branching ratio
#--------------------------
def get_BR(gSM, gDM, mV, mDM):
    return get_width_inv(gDM,mV,mDM)/get_total_width(gSM, gDM, mV, mDM)



# gDM from total width
#---------------------
def get_gDM_from_width(width, gSM, mV, mDM):
    w_vis = get_width_vis(gSM,mV)
    gDM2  = (width-w_vis) / PhiInv(mDM,mV)
    return sqrt(gDM2)


# gDM from invisble BR
#---------------------
def get_gDM_from_BR(BR, gSM, mV, mDM):
    gDM2 = gSM**2 * (BR/(1-BR)) * PhiVis(mt,mV) / PhiInv(mDM,mV)
    return sqrt(gDM2)



# gSM from invisble BR and total width
#-------------------------------------
def get_gSM_from_BRwidth(BR, width, mV, mDM):
    gSM2 = width / PhiVis(mt ,mV) * (1-BR)
    return sqrt(gSM2)

get_gSM_from_BRwidth=np.vectorize(get_gSM_from_BRwidth)

# gDM from invisble BR and total width
#-------------------------------------
def get_gDM_from_BRwidth(BR, width, mV, mDM):
    gDM2 = width / PhiInv(mDM,mV) * BR
    return sqrt(gDM2)

get_gDM_from_BRwidth=np.vectorize(get_gDM_from_BRwidth)

